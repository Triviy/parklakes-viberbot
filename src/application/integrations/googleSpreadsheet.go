package integrations

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// GoogleSpreadsheet used for Google Spreadsheets API
type GoogleSpreadsheet struct {
	service       *sheets.Service
	spreadsheetID string
}

// NewGoogleSpreadsheet returns new GoogleSpreadsheet
func NewGoogleSpreadsheet(ctx context.Context, apiKey string, spreadsheetID string) (*GoogleSpreadsheet, error) {
	srv, err := sheets.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, errors.Wrap(err, "creation of spreadsheet service failed")
	}
	return &GoogleSpreadsheet{srv, spreadsheetID}, nil
}

// ReadSpreadsheetRange reads range of Google Spreadsheet values
func (gs GoogleSpreadsheet) ReadSpreadsheetRange(readRange string) ([][]interface{}, error) {
	resp, err := gs.service.Spreadsheets.Values.Get(gs.spreadsheetID, readRange).Do()
	if err != nil {
		return nil, errors.Wrap(err, "getting data from spreadsheet failed")
	}
	return resp.Values, nil
}
