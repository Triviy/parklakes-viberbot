package main

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type apiConfig struct {
	AppBaseURL         string `yaml:"appBaseURL" env:"WEBSITE_HOSTNAME" env-default:"http://localhost"`
	APIKey             string `yaml:"apiKey" env:"API_KEY" env-required:"true"`
	InstrumentationKey string `yaml:"instrumentationKey" env:"APPINSIGHTS_INSTRUMENTATION_KEY" env-required:"true"`
}

func getAPIConfig() (*apiConfig, error) {
	var apiConfig apiConfig
	cfgFile := "config.yml"
	if fileExists(cfgFile) {
		err := cleanenv.ReadConfig("config.yml", &apiConfig)
		if err != nil {
			return nil, err
		}
	}
	err := cleanenv.ReadEnv(&apiConfig)
	if err != nil {
		return nil, err
	}
	return &apiConfig, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
