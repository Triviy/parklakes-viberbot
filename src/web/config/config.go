package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type apiConfig struct {
	Database struct {
		ConnectionString string `yaml:"connectionString" env:"DB_CONNECTION_STRING" env-required:"true"`
	} `yaml:"database"`
	Viber struct {
		APIKey     string `yaml:"apiKey" env:"VIBER_API_KEY" env-required:"true"`
		WebhookURL string `yaml:"webhookURL" env:"VIBER_WEBHOOK_URL" env-required:"true"`
		BaseURL    string `yaml:"baseURL" env:"VIBER_BASE_URL" env-required:"true"`
	} `yaml:"viber"`
	SheetsAPI struct {
		SpreadsheetID string `yaml:"spreadsheetID" env:"SHEETS_SPREADSHEET_ID" env-required:"true"`
		APIKey        string `yaml:"apiKey" env:"SHEETS_API_KEY" env-required:"true"`
	} `yaml:"sheetsAPI"`
}

// APIConfig ...
type APIConfig struct {
	cfg *apiConfig
}

// NewAPIConfig initalizes configuration for application
func NewAPIConfig() (apiCfg *APIConfig, err error) {
	cfgFile := "config.yml"
	if fileExists(cfgFile) {
		err = cleanenv.ReadConfig("config.yml", &apiCfg.cfg)
	}
	err = cleanenv.ReadEnv(&apiCfg.cfg)
	return
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetDBConnectionString returns database connection string from configuration
func (c APIConfig) GetDBConnectionString() string {
	return c.cfg.Database.ConnectionString
}

// GetViberAPIKey returns Viber AuthToken
func (c APIConfig) GetViberAPIKey() string {
	return c.cfg.Viber.APIKey
}

// GetViberWebhookURL returns Viber Webhook URL
func (c APIConfig) GetViberWebhookURL() string {
	return c.cfg.Viber.WebhookURL
}

// GetViberBaseURL returns Viber chatbot API URL
func (c APIConfig) GetViberBaseURL() string {
	return c.cfg.Viber.BaseURL
}

// GetSheetsAPISpreadsheetID returns ID of car owners Google SpreadSheet
func (c APIConfig) GetSheetsAPISpreadsheetID() string {
	return c.cfg.SheetsAPI.SpreadsheetID
}

// GetSheetsAPIKey returns APIKey with scoped access to Google SpreadSheet
func (c APIConfig) GetSheetsAPIKey() string {
	return c.cfg.SheetsAPI.APIKey
}
