package main

import (
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

func newAppinsightsClient(key string) appinsights.TelemetryClient {
	telemetryConfig := appinsights.NewTelemetryConfiguration(key)
	telemetryConfig.MaxBatchSize = 5
	telemetryConfig.MaxBatchInterval = time.Second
	client := appinsights.NewTelemetryClientFromConfig(telemetryConfig)
	client.Context().Tags.Cloud().SetRole("parklakes-migration-func")
	return client
}

func closeAppinsightsClient(ac appinsights.TelemetryClient) {
	select {
	case <-ac.Channel().Close(3 * time.Second):
	case <-time.After(10 * time.Second):
	}
}
