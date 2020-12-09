package github.com/triviy/parklakes-viberbot/handlers

// MigrateCarOwners runs migration from Google SpreadSheet to database
func MigrateCarOwners() {
	spreadsheetID := GetSheetsAPISpreadsheetID()
	carOwners := GetCarOwners(spreadsheetID)
	MigrateCardOwners(carOwners)
	GetCardOwner("CB4498CC")
}
