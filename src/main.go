package main

import (
	"context"
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/triviy/parklakes-viberbot/infrastructure/logger"
	"github.com/triviy/parklakes-viberbot/web"
	"github.com/triviy/parklakes-viberbot/web/config"
	"github.com/triviy/parklakes-viberbot/web/handlers"
	"github.com/triviy/parklakes-viberbot/web/middlewares"
)

func main() {

	logger.Info("Starting main")

	ctx := context.Background()
	cfg, err := config.NewAPIConfig()
	if err != nil {
		logger.Fatal(err)
		return
	}

	h, err := handlers.InitializeHandlers(ctx, cfg)
	if err != nil {
		logger.Fatal(err)
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

	port := fmt.Sprintf(":%s", cfg.GetAppPort())
	web.GracefulShutdown(e, e.Start(port))
}
