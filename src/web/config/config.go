package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

type apiConfig struct {
	AppPort    string `yaml:"appPort" env:"PORT" env-default:"8081"`
	AppBaseURL string `yaml:"appBaseURL" env:"WEBSITE_HOSTNAME" env-default:"http://localhost"`
	APIKey     string `yaml:"apiKey" env:"API_KEY" env-required:"true"`
	Database   struct {
		ConnectionString string `yaml:"connectionString" env:"DB_CONNECTION_STRING" env-required:"true"`
	} `yaml:"database"`
	Viber struct {
		APIKey  string `yaml:"apiKey" env:"VIBER_API_KEY" env-required:"true"`
		BaseURL string `yaml:"baseURL" env:"VIBER_BASE_URL" env-required:"true"`
	} `yaml:"viber"`
	SheetsAPI struct {
		SpreadsheetID string `yaml:"spreadsheetID" env:"SHEETS_SPREADSHEET_ID" env-required:"true"`
		APIKey        string `yaml:"apiKey" env:"SHEETS_API_KEY" env-required:"true"`
	} `yaml:"sheetsAPI"`
	ComputerVision struct {
		APIKey string `yaml:"apiKey" env:"COMPUTER_VISION_API_KEY" env-required:"true"`
		APIUrl string `yaml:"apiURL" env:"COMPUTER_VISION_API_URL" env-required:"true"`
	} `yaml:"computerVision"`
	AppInsights struct {
		InstrumentationKey string `yaml:"instrumentationKey" env:"APPINSIGHTS_INSTRUMENTATION_KEY" env-required:"true"`
	} `yaml:"appInsights"`
}

// APIConfig ...
type APIConfig struct {
	cfg *apiConfig
}

// NewAPIConfig initalizes configuration for application
func NewAPIConfig() (*APIConfig, error) {
	var apiConfig apiConfig
	cfgFile := "config.yml"
	if fileExists(cfgFile) {
		err := cleanenv.ReadConfig("config.yml", &apiConfig)
		if err != nil {
			return nil, errors.Wrap(err, "creating config failed")
		}
	}
	err := cleanenv.ReadEnv(&apiConfig)
	if err != nil {
		return nil, errors.Wrap(err, "creating config failed")
	}
	return &APIConfig{&apiConfig}, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetAppPort returns application running port
func (c APIConfig) GetAppPort() string {
	return c.cfg.AppPort
}

// GetAPIKey returns API key required for internal use API
func (c APIConfig) GetAPIKey() string {
	return c.cfg.APIKey
}

// GetAppBaseURL returns application base URL
func (c APIConfig) GetAppBaseURL() string {
	return c.cfg.AppBaseURL
}

// GetDBConnectionString returns database connection string from configuration
func (c APIConfig) GetDBConnectionString() string {
	return c.cfg.Database.ConnectionString
}

// GetViberAPIKey returns Viber AuthToken
func (c APIConfig) GetViberAPIKey() string {
	return c.cfg.Viber.APIKey
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

// GetComputerVisionAPIKey returns APIKey for Computer Vision API
func (c APIConfig) GetComputerVisionAPIKey() string {
	return c.cfg.ComputerVision.APIKey
}

// GetComputerVisionAPIUrl returns URL for Computer Vision API
func (c APIConfig) GetComputerVisionAPIUrl() string {
	return c.cfg.ComputerVision.APIUrl
}

// GetAppInsightsInstrumentationKey returns InstrumentationKey for AppInsights
func (c APIConfig) GetAppInsightsInstrumentationKey() string {
	return c.cfg.AppInsights.InstrumentationKey
}
