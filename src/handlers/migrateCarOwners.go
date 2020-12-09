package handlers

import (
	"github.com/triviy/parklakes-viberbot/config"
	"github.com/triviy/parklakes-viberbot/stores"
)

// MigrateCarOwners runs migration from Google SpreadSheet to database
func MigrateCarOwners() {
	spreadsheetID := config.GetSheetsAPISpreadsheetID()
	carOwners := stores.GetCarOwners(spreadsheetID)
	stores.MigrateCardOwners(carOwners)
	stores.GetCardOwner("CB4498CC")
}
