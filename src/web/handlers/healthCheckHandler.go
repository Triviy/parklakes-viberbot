package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/triviy/parklakes-viberbot/infrastructure/persistance"
)

type healthResponse struct {
	Status  string          `json:"status"`
	Details map[string]bool `json:"details"`
}

const (
	ok   = "OK"
	fail = "FAIL"
)

// HealthCheckHandler handles migration requests
type HealthCheckHandler struct {
	datastore *persistance.MongoDatastore
}

// NewHealthCheckHandler creates new handler instance
func NewHealthCheckHandler(ds *persistance.MongoDatastore) *HealthCheckHandler {
	return &HealthCheckHandler{ds}
}

// Handle runs migration from Google SpreadSheet to database
func (h HealthCheckHandler) Handle(c echo.Context) error {
	status := http.StatusOK
	m := healthResponse{Status: ok}

	dbCheck := h.datastore.Ping()
	if dbCheck != nil {
		status = http.StatusInternalServerError
		m.Status = fail
		m.Details["db"] = false
	}

	return c.JSON(status, m)
}
