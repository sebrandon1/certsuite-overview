package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBUser      string
	DBPassword  string
	DBURL       string
	DBPort      string
	ClientID    string
	APISecret   string
	BearerToken string
	Namespace   string
	Repository  string
}

var AppConfig Config

// Initialize Viper and load configuration
func LoadConfig() {
	// Configure Viper to read from environment variables
	viper.AutomaticEnv()

	// Load the configuration into the AppConfig struct
	AppConfig = Config{
		DBUser:      GetConfigValue("DB_USER"),
		DBPassword:  GetConfigValue("DB_PASSWORD"),
		DBURL:       GetConfigValue("DB_URL"),
		DBPort:      GetConfigValue("DB_PORT"),
		ClientID:    GetConfigValue("CLIENTID"),
		APISecret:   GetConfigValue("APISECRET"),
		BearerToken: GetConfigValue("BEARERTOKEN"),
		Namespace:   GetConfigValue("NAMESPACE"),
		Repository:  GetConfigValue("REPOSITORY"),
	}
}

// Helper function to get a configuration value by key
func GetConfigValue(key string) string {
	value := viper.GetString(key)
	if value == "" {
		log.Fatalf("Configuration key %s is missing", key)
	}
	return value
}