package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/spf13/viper"

	_ "github.com/spf13/viper/remote" // imported blank to enable local AND remote config read
)

// Configuration represents the configuration object
type Configuration struct {
	LogLevel string            `json:"-"`
	S3       S3                `json:"-"`
	OIDC     OidcConfiguration `json:"-"`
}

var configuration Configuration
var consulURL string

// InitConfig is called on service startup to bring in configuration variables
func InitConfig() Configuration {

	// If not release mode, use default consul url. Comment out this block if you want to use local config.yaml.
	ginMode := os.Getenv("GIN_MODE")
	if ginMode != "release" {
		os.Setenv("ConfigurationUrl", "localhost:8500")
	}

	viper.BindEnv("CONFIGURL", "ConfigurationUrl")
	consulURL = viper.GetString("CONFIGURL")

	if consulURL == "" {
		replacer := strings.NewReplacer(".", "__")
		viper.SetEnvKeyReplacer(replacer)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}

		err := viper.Unmarshal(&configuration)
		if err != nil {
			log.Fatalf("Unable to decode into struct, %v", err)
		}
		return configuration
	}

	fmt.Println("Connecting to consul at " + consulURL)
	// Else, read from consul
	viper.AddRemoteProvider("consul", consulURL, "tcp-auditor/appsettings.json")
	viper.SetConfigType("json") // Set to JSON for consul
	viper.AutomaticEnv()
	if err := viper.ReadRemoteConfig(); err != nil {
		log.Fatalf("Error reading config from consul, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return configuration
}

// GetConfig returns the configuration to be read by the service
func GetConfig() Configuration {
	return configuration
}
