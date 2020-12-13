package main

import (
	"context"
	"fmt"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/triviy/parklakes-viberbot/web/config"
	"github.com/triviy/parklakes-viberbot/web/handlers"
	"github.com/triviy/parklakes-viberbot/web/middlewares"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewAPIConfig()
	if err != nil {
		log.Fatal(err)
	}

	h, err := handlers.InitializeHandlers(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	apiKeyAuth := middlewares.GetAPIKeyAuthMiddleware(cfg.GetAPIKey())
	e.POST("/api/v1/car-owners/migrate", h.MigrateCarOwnersHandler.Handle, apiKeyAuth)
	// e.POST("/api/v1/viber/set-webhook", handlers.SetWebhook)
	// e.POST("/api/v1/viber/callback", handlers.SendMessage)
	port := fmt.Sprintf(":%s", cfg.GetAppPort())
	e.Logger.Fatal(e.Start(port))
}
