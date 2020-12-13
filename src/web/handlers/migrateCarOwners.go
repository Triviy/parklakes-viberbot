package handlers

import (
	"github.com/labstack/echo"
	"github.com/triviy/parklakes-viberbot/config"
	"github.com/triviy/parklakes-viberbot/stores"
)

// MigrateCarOwners runs migration from Google SpreadSheet to database
func MigrateCarOwners(c echo.Context) error {
	spreadsheetID := config.GetSheetsAPISpreadsheetID()
	carOwners := stores.GetCarOwners(spreadsheetID)
	stores.MigrateCardOwners(carOwners)
}
