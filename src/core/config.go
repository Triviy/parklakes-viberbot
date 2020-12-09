package github.com/triviy/parklakes-viberbot/core

import (
	"github.com/ilyakaznacheev/cleanenv"
)

var cfg apiConfig

type apiConfig struct {
	Database struct {
		ConnectionString string `yaml:"connectionString" envconfig:"DB_CONNECTION_STRING"`
	} `yaml:"database"`
	Viber struct {
		APIKey     string `yaml:"apiKey" envconfig:"VIBER_API_KEY"`
		WebhookURL string `yaml:"webhookURL" envconfig:"VIBER_WEBHOOK_URL"`
		BaseURL    string `yaml:"baseURL" envconfig:"VIBER_BASE_URL"`
	} `yaml:"viber"`
	SheetsAPI struct {
		SpreadsheetID   string `yaml:"spreadsheetID" envconfig:"SHEETS_SPREADSHEET_ID"`
		CredentialsJSON string `yaml:"credentialsJSON" envconfig:"SHEETS_API_CREDENTIALS_JSON"`
		TokenJSON       string `yaml:"tokenJSON" envconfig:"SHEETS_API_TOKEN_JSON"`
	} `yaml:"sheetsAPI"`
}

// InitalizeAPIConfig initalizes configuration for application
func InitalizeAPIConfig() error {
	return cleanenv.ReadConfig("../config.yml", &cfg)
}

// GetDBConnectionString returns database connection string from configuration
func GetDBConnectionString() string {
	return cfg.Database.ConnectionString
}

// GetViberAPIKey returns Viber AuthToken
func GetViberAPIKey() string {
	return cfg.Viber.APIKey
}

// GetViberWebhookURL returns Viber Webhook URL
func GetViberWebhookURL() string {
	return cfg.Viber.WebhookURL
}

// GetViberBaseURL returns Viber chatbot API URL
func GetViberBaseURL() string {
	return cfg.Viber.BaseURL
}

// GetSheetsAPISpreadsheetID returns ID of car owners Google SpreadSheet
func GetSheetsAPISpreadsheetID() string {
	return cfg.SheetsAPI.SpreadsheetID
}

// GetSheetsAPICredentialsJSON returns JSON credentials for authentication to Google SpreadSheet API
func GetSheetsAPICredentialsJSON() string {
	return cfg.SheetsAPI.CredentialsJSON
}

// GetSheetsAPITokenJSON returns JSON OAuth tokens obtained with GetSheetsAPICredentialsJSON
func GetSheetsAPITokenJSON() string {
	return cfg.SheetsAPI.TokenJSON
}
