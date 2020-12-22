package logging

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
	"github.com/sirupsen/logrus"
)

var defaultLevels = []logrus.Level{
	logrus.PanicLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
}

var levelMap = map[logrus.Level]contracts.SeverityLevel{
	logrus.PanicLevel: appinsights.Critical,
	logrus.FatalLevel: appinsights.Critical,
	logrus.ErrorLevel: appinsights.Error,
	logrus.WarnLevel:  appinsights.Warning,
	logrus.InfoLevel:  appinsights.Information,
}

// AppInsightsHook is a logrus hook for Application Insights
type AppInsightsHook struct {
	client appinsights.TelemetryClient

	async        bool
	levels       []logrus.Level
	ignoreFields map[string]struct{}
	filters      map[string]func(interface{}) interface{}
}

// NewAppInsightsHook returns an initialised logrus hook for Application Insights
func NewAppInsightsHook(name string, instrumentationKey string) *AppInsightsHook {
	tc := appinsights.NewTelemetryConfiguration(instrumentationKey)
	tc.MaxBatchSize = 8192
	tc.MaxBatchInterval = 2 * time.Second

	c := appinsights.NewTelemetryClientFromConfig(tc)

	if name != "" {
		c.Context().Tags.Cloud().SetRole(name)
	}
	return &AppInsightsHook{
		client:       c,
		levels:       defaultLevels,
		ignoreFields: make(map[string]struct{}),
		filters:      make(map[string]func(interface{}) interface{}),
	}
}

// Levels returns logging level to fire this hook.
func (hook *AppInsightsHook) Levels() []logrus.Level {
	return hook.levels
}

// SetLevels sets logging level to fire this hook.
func (hook *AppInsightsHook) SetLevels(levels []logrus.Level) {
	hook.levels = levels
}

// SetAsync sets async flag for sending logs asynchronously.
// If use this true, Fire() does not return error.
func (hook *AppInsightsHook) SetAsync(async bool) {
	hook.async = async
}

// AddIgnore adds field name to ignore.
func (hook *AppInsightsHook) AddIgnore(name string) {
	hook.ignoreFields[name] = struct{}{}
}

// AddFilter adds a custom filter function.
func (hook *AppInsightsHook) AddFilter(name string, fn func(interface{}) interface{}) {
	hook.filters[name] = fn
}

// Fire is invoked by logrus and sends log data to Application Insights.
func (hook *AppInsightsHook) Fire(entry *logrus.Entry) error {
	if !hook.async {
		return hook.fire(entry)
	}
	// async - fire and forget
	go hook.fire(entry)
	return nil
}

func (hook *AppInsightsHook) fire(entry *logrus.Entry) error {
	trace, err := hook.buildTrace(entry)
	if err != nil {
		return err
	}
	hook.client.Track(trace)
	return nil
}

func (hook *AppInsightsHook) buildTrace(entry *logrus.Entry) (*appinsights.TraceTelemetry, error) {
	// Add the message as a field if it isn't already
	if _, ok := entry.Data["message"]; !ok {
		entry.Data["message"] = entry.Message
	}

	level := levelMap[entry.Level]
	trace := appinsights.NewTraceTelemetry(entry.Message, level)
	if trace == nil {
		return nil, fmt.Errorf("Could not create telemetry trace with entry %+v", entry)
	}
	for k, v := range entry.Data {
		if _, ok := hook.ignoreFields[k]; ok {
			continue
		}
		if fn, ok := hook.filters[k]; ok {
			v = fn(v) // apply custom filter
		} else {
			v = formatData(v) // use default formatter
		}
		vStr := fmt.Sprintf("%v", v)
		trace.Properties[k] = vStr
	}
	trace.Properties["source_level"] = entry.Level.String()
	trace.Timestamp = entry.Time
	return trace, nil
}

// formatData returns value as a suitable format.
func formatData(value interface{}) (formatted interface{}) {
	switch value := value.(type) {
	case json.Marshaler:
		return value
	case error:
		return value.Error()
	case fmt.Stringer:
		return value.String()
	default:
		return value
	}
}

func stringPtr(str string) *string {
	return &str
}
