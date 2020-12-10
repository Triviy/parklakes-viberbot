package main

import (
	"log"

	"github.com/labstack/echo"
	"github.com/triviy/parklakes-viberbot/config"
	"github.com/triviy/parklakes-viberbot/handlers"
)

func main() {
	err := config.InitalizeAPIConfig()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.POST("/api/v1/car-owners/migrate", func(c echo.Context) error {
		handlers.MigrateCarOwners()
		return nil
	})
	e.POST("/api/v1/viber/set-webhook", func(c echo.Context) error {
		handlers.SetWebhook()
		return nil
	})
	e.POST("/api/v1/viber/callback", func(c echo.Context) error {
		c.Request()
		handlers.SendMessage("asdasd", "asdasd")
		return nil
	})
	e.Logger.Fatal(e.Start(":8081"))
}
