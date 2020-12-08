package main

import (
	"log"
)

func main() {
	err := InitalizeAPIConfig()
	if err != nil {
		log.Fatal(err)
	}

	//MigrateCarOwners()
	//SetWebhook()
	SendMessage("asdasd", "asdasd")

	// e := echo.New()
	// e.POST("/api/v1/car-owners/migrate", func(c echo.Context) error {
	// 	MigrateCarOwners()
	// 	return nil
	// })
	// e.POST("/api/v1/viber/set-webhook", func(c echo.Context) ecdrror {
	// 	SetWebhook()
	// 	return nil
	// })
	// e.POST("/api/v1/viber/callback", func(c echo.Context) error {
	// 	c.Request()
	// 	SendMessage("asdasd", "asdasd")
	// 	return nil
	// })
	// e.Logger.Fatal(e.Start(":8081"))
}
