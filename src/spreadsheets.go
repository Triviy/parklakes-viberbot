package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// GetCarOwners gets car owners data from Google SpreadSheet
func GetCarOwners(spreadsheetID string) map[string]CarOwner {
	srv := getSpreadsheetService()
	readRange := "A2:E"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	log.Printf("Got %v records:\n", len(resp.Values))

	carOwners := make(map[string]CarOwner)
	for _, row := range resp.Values {
		carOwner := CreateCarOwnerFromRecord(row)
		if carOwner != nil {
			carOwners[carOwner.ID] = *carOwner
		}
	}
	log.Printf("Got %v unique records\n", len(carOwners))
	return carOwners
}

func getSpreadsheetService() *sheets.Service {
	b := []byte(GetSheetsAPICredentialsJSON())
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	return srv
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	token := &oauth2.Token{}
	err := json.Unmarshal([]byte(GetSheetsAPITokenJSON()), token)
	if err != nil {
		token = getTokenFromWeb(config)
		saveToken("token.json", token)
	}
	return config.Client(context.Background(), token)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	log.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
