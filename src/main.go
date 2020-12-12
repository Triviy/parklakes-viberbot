package main

import (
	"log"

	"github.com/labstack/echo"
	"github.com/triviy/parklakes-viberbot/config"
	"github.com/triviy/parklakes-viberbot/handlers"
	"github.com/triviy/parklakes-viberbot/web/middlewares"
)

func main() {
	cfg, err := config.NewAPIConfig()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middlewares.ExtendWithAppContext())
	e.POST("/api/v1/car-owners/migrate", handlers.MigrateCarOwners)
	e.POST("/api/v1/viber/set-webhook", handlers.SetWebhook)
	e.POST("/api/v1/viber/callback", handlers.SendMessage)
	e.Logger.Fatal(e.Start(":8081"))
}
