package logging

import (
	"os"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/sirupsen/logrus"
)

// InitLog initializes logging
func InitLog(instrumentationKey string) {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	logrus.SetOutput(os.Stdout)
	hook := NewAppInsightsHook("parklakes-viberbot", &appinsights.TelemetryConfiguration{
		InstrumentationKey: instrumentationKey,
		MaxBatchSize:       8192,
		MaxBatchInterval:   2 * time.Second,
	})
	logrus.AddHook(hook)
}
