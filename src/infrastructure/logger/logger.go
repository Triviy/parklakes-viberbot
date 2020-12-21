package logger

// TODO: add request_id in log entry, add interface, extend echo.Context
// TODO: leave only appinsights or combine with logrus as decorator
// TODO: wrap Repo and jsonHttp with TrackDependency

import (
	"fmt"
	"os"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/sirupsen/logrus"
)

const (
	detailProp = "details"
)

type requestLogEntry struct {
	Method        string `json:"method"`
	URI           string `json:"uri"`
	Path          string `json:"path"`
	RemoteIP      string `json:"remoteIp"`
	Host          string `json:"host"`
	Protocol      string `json:"protocol"`
	Referer       string `json:"referer"`
	UserAgent     string `json:"userAgent"`
	RequestID     string `json:"requestId"`
	Body          string `json:"body"`
	ContentLength int64  `json:"contentLength"`
}

type responseLogEntry struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
	Method     string `json:"method"`
	URI        string `json:"uri"`
	Path       string `json:"path"`
	RequestID  string `json:"requestId"`
	Latency    string `json:"latency"`
	StackTrace string `json:"stackTrace"`
}

var client appinsights.TelemetryClient

// InitializeLogger sets up logger instance
func InitializeLogger(ik string) {
	telemetryConfig := appinsights.NewTelemetryConfiguration(ik)
	telemetryConfig.MaxBatchSize = 8192
	telemetryConfig.MaxBatchInterval = 2 * time.Second
	client = appinsights.NewTelemetryClientFromConfig(telemetryConfig)
	client.Context().Tags.Cloud().SetRole("parklakes-viberbot")

	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	logrus.SetOutput(os.Stdout)
}

// Info logs a message at level Info on the standard logger
func Info(args ...interface{}) {
	logrus.Info(args...)
	client.TrackTrace(fmt.Sprint(args...), appinsights.Information)
}

// InfoDetailed logs a message at level Info on the standard logger with additional data
func InfoDetailed(detail interface{}, args ...interface{}) {
	logrus.WithField(detailProp, detail).Info(args...)

	trace := appinsights.NewTraceTelemetry(fmt.Sprint(args...), appinsights.Information)
	trace.Properties[detailProp] = fmt.Sprint(detail)
	client.Track(trace)
}

// Infof logs a message at level Info on the standard logger
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args)
	client.TrackTrace(fmt.Sprintf(format, args...), appinsights.Information)
}

// InfofDetailed logs a message at level Info on the standard logger with additional data
func InfofDetailed(detail interface{}, format string, args ...interface{}) {
	logrus.WithField(detailProp, detail).Infof(format, args)

	trace := appinsights.NewTraceTelemetry(fmt.Sprintf(format, args...), appinsights.Information)
	trace.Properties[detailProp] = fmt.Sprint(detail)
	client.Track(trace)
}

// Warn logs a message at level Warn on the standard logger
func Warn(args ...interface{}) {
	logrus.Info(args...)
	client.TrackTrace(fmt.Sprint(args...), appinsights.Warning)
}

// Warnf logs a message at level Warn on the standard logger
func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args)
	client.TrackTrace(fmt.Sprintf(format, args...), appinsights.Warning)
}

// Error logs a message at level Error on the standard logger
func Error(args ...interface{}) {
	logrus.Error(args...)
	client.TrackTrace(fmt.Sprint(args...), appinsights.Error)
}

// Errorf logs a message at level Error on the standard logger
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args)
	client.TrackTrace(fmt.Sprintf(format, args...), appinsights.Error)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	defer appinsights.TrackPanic(client, false)
	logrus.Fatal(args)
}

// Request registers application requests
func Request(r requestLogEntry) {
	InfoDetailed(r, "-- Start request")
}

// Response registers application responses
func Response(r responseLogEntry) {
	logrus.WithField(detailProp, r).Info("-- End request")

	var code string
	if r.Error == "" {
		code = fmt.Sprint(r.StatusCode)
	} else {
		code = fmt.Sprintf("%v - %s", r.StatusCode, r.Error)
	}
	d, _ := time.ParseDuration(r.Latency)
	client.TrackRequest(r.Method, r.URI, d, code)
}

// Close finishes submitting of telemetry and closes channel
func Close() {
	select {
	case <-client.Channel().Close(10 * time.Second):
		// Ten second timeout for retries.

		// If we got here, then all telemetry was submitted
		// successfully, and we can proceed to exiting.
	case <-time.After(30 * time.Second):
		// Thirty second absolute timeout.  This covers any
		// previous telemetry submission that may not have
		// completed before Close was called.

		// There are a number of reasons we could have
		// reached here.  We gave it a go, but telemetry
		// submission failed somewhere.  Perhaps old events
		// were still retrying, or perhaps we're throttled.
		// Either way, we don't want to wait around for it
		// to complete, so let's just exit.
	}
}
