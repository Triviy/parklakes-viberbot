package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/sirupsen/logrus"
)

// InitLog initializes logging
func InitLog(instrumentationKey string, enableTracingDiag bool) {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	logrus.SetOutput(os.Stdout)
	hook := NewAppInsightsHook("parklakes-viberbot", instrumentationKey)
	logrus.AddHook(hook)
	if enableTracingDiag {
		appinsights.NewDiagnosticsMessageListener(func(msg string) error {
			fmt.Printf("[%s] %s\n", time.Now().Format(time.UnixDate), msg)
			return nil
		})
	}
}
