package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

func main() {
	cfg, err := getAPIConfig()
	if err != nil {
		log.Fatal(err)
	}

	ac := newAppinsightsClient(cfg.InstrumentationKey)
	ac.TrackTrace("migration started", appinsights.Information)

	url := fmt.Sprintf("%s/api/v1/car-owners/migrate", cfg.AppBaseURL)
	r, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		ac.TrackTrace(fmt.Sprintf("creating request for migrate API failed: %s", err.Error()), appinsights.Error)
		log.Fatal(err)
	}
	r.Header.Add("X-API-KEY", cfg.APIKey)

	c := &http.Client{}
	res, err := c.Do(r)
	if err != nil {
		ac.TrackTrace(fmt.Sprintf("calling migrate API failed: %s", err.Error()), appinsights.Error)
		log.Fatal(err)
	}
	ac.TrackTrace(fmt.Sprintf("migration ended with status %s", res.Status), appinsights.Information)

	closeAppinsightsClient(ac)
}
