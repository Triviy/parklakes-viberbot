package web

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
	"github.com/triviy/parklakes-viberbot/infrastructure/logger"
	"github.com/triviy/parklakes-viberbot/infrastructure/persistance"
)

// GracefulShutdown Wait for interrupt signal
// to gracefully shutdown the server with a timeout of 5 seconds.
func GracefulShutdown(e *echo.Echo, err error) {
	logger.Error(err)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if datastore := persistance.GetDatastore(); datastore != nil {
		err := datastore.Disconnect()
		if err != nil {
			logger.Error(err)
		}
	}
	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}
}
