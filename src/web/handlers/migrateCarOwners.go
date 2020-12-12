package handlers

import (
	"github.com/labstack/echo"
	"github.com/triviy/parklakes-viberbot/config"
	"github.com/triviy/parklakes-viberbot/stores"
)

// MigrateCarOwners runs migration from Google SpreadSheet to database
func MigrateCarOwners(c echo.Context) error {
	// kts, err := core.GetKyivTimeString()
	// if err != nil {
	// 	return errors.Wrap(err, "Getting current Kyiv time string failed")
	// }
	spreadsheetID := config.GetSheetsAPISpreadsheetID()
	carOwners := stores.GetCarOwners(spreadsheetID)
	stores.MigrateCardOwners(carOwners)
	stores.GetCardOwner("CB4498CC")
}
