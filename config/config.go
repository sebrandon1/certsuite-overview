package config

import (
	"encoding/json"
	"os"
)

// Config struct to hold the configuration values
type Config struct {
	BearerToken string `json:"bearerToken"`
	ClientID    string `json:"clientID"`
	APISecret   string `json:"apiSecret"`
	Namespace   string `json:"namespace"`
	Repository  string `json:"repository"`
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filePath string) (*Config, error) {
	// Open the config file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the JSON file into the Config struct
	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
