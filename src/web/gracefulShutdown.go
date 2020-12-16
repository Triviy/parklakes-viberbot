package web

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/triviy/parklakes-viberbot/infrastructure/persistance"
)

// GracefulShutdown Wait for interrupt signal
// to gracefully shutdown the server with a timeout of 5 seconds.
func GracefulShutdown(e *echo.Echo, err error) {
	log.Info("we are in GracefulShutdown")
	log.Error(err)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Info("we are closing datastore")
	if datastore := persistance.GetDatastore(); datastore != nil {
		err := datastore.Disconnect()
		if err != nil {
			log.Error(err)
		}
	}
	log.Info("we are at Shutdown")
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
