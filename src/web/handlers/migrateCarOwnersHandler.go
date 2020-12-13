package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/triviy/parklakes-viberbot/application/commands"
)

// MigrateCarOwnersHandler handles migration requests
type MigrateCarOwnersHandler struct {
	cmd *commands.MigrateCarOwnersCmd
}

// NewMigrateCarOwnersHandler creates new handler instance
func NewMigrateCarOwnersHandler(cmd *commands.MigrateCarOwnersCmd) *MigrateCarOwnersHandler {
	return &MigrateCarOwnersHandler{cmd}
}

// Handle runs migration from Google SpreadSheet to database
func (h MigrateCarOwnersHandler) Handle(c echo.Context) error {
	err := h.cmd.Migrate()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, createOkResponse())
}
