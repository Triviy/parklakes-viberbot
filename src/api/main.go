package main

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/triviy/parklakes-viberbot/infrastructure/logging"
	"github.com/triviy/parklakes-viberbot/web"
	"github.com/triviy/parklakes-viberbot/web/config"
	"github.com/triviy/parklakes-viberbot/web/handlers"
	"github.com/triviy/parklakes-viberbot/web/middlewares"
)

var (
	cfg *config.APIConfig
)

func init() {
	c, err := config.NewAPIConfig()
	if err != nil {
		log.Fatal(err)
		return
	}
	logging.InitLog(c.GetAppInsightsInstrumentationKey(), c.GetAppInsightsEnableTracingDiag())
	cfg = c
}

func main() {
	log.Info("Starting main")

	h, err := handlers.InitializeHandlers(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middlewares.CustomLogger())

	apiKeyAuth := middlewares.APIKeyAuth(cfg.GetAPIKey())

	e.GET("/health", h.HealthCheckHandler.Handle)
	e.POST("/api/v1/car-owners/migrate", h.MigrateCarOwnersHandler.Handle, apiKeyAuth)
	e.POST("/api/v1/viber/set-webhook", h.SetWebhookHandler.Handle, apiKeyAuth)
	e.POST("/api/v1/viber/callback", h.CallbackHandler.Handle, middlewares.ViberHashCheck(cfg.GetViberAPIKey()))
	e.File("/", fmt.Sprintf("%s/index.html", cfg.GetAssetsPath()))

	port := fmt.Sprintf(":%s", cfg.GetAppPort())
	web.GracefulShutdown(e, e.Start(port))
}
