package google

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Spreadsheet used for Google Spreadsheets API
type Spreadsheet struct {
	service       *sheets.Service
	spreadsheetID string
}

// NewSpreadsheet returns new Spreadsheet
func NewSpreadsheet(ctx context.Context, apiKey string, spreadsheetID string) (*Spreadsheet, error) {
	srv, err := sheets.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, errors.Wrap(err, "creation of spreadsheet service failed")
	}
	return &Spreadsheet{srv, spreadsheetID}, nil
}

// ReadSpreadsheetRange reads range of Google Spreadsheet values
func (gs Spreadsheet) ReadSpreadsheetRange(readRange string) ([][]interface{}, error) {
	resp, err := gs.service.Spreadsheets.Values.Get(gs.spreadsheetID, readRange).Do()
	if err != nil {
		return nil, errors.Wrap(err, "getting data from spreadsheet failed")
	}
	return resp.Values, nil
}
