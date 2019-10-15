package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config represent the application configuration
type Config struct {
	ServerName         string `json:"server_name"`
	KerioStorePath     string `json:"kerioStore_path"`
	Sender             string `json:"sender"`
	Recipient          string `json:"recipient"`
	SubjectT           string `json:"subject_t"`
	HTMLBodyT          string `json:"html_body"`
	TextBodyT          string `json:"text_body"`
	CharSet            string `json:"char_set"`
	QueuePath          string `json:"q_path"`
	QueueWarnThreshold int    `json:"q_warn_threshold"`
}

// New Inits a new instance of config.
func New() (*Config, error) {
	// initialize the configuration variables, either from config file or using defaults
	viper.SetConfigName("kerio-checks-config") // name of config file
	viper.AddConfigPath("/opt/kerio-check/")   // path to look for the config file in
	viper.AddConfigPath(".")                   // optionally look for config in the working directory
	viper.SetDefault("ServerName", "Kerio-TEST")
	viper.SetDefault("KerioStorePath", "/mnt/datastore01/kerio/mailserver/store/")
	viper.SetDefault("Sender", "noreply-ses@aws.example.tld")
	viper.SetDefault("Recipient", "support-alert-crew@example.tld")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("No config file found, proceeding with defaults, and creating a new config file in the working directory")
			viper.SetConfigType("json")
			err = viper.WriteConfigAs("./kerio-checks-config.json")
			if err != nil {
				return nil, fmt.Errorf("Fatal error writing config file: %s", err)
			}
		} else {
			// Config file was found but another error was produced
			return nil, fmt.Errorf("Fatal error config file: %s", err)
		}
	}
	config := &Config{}
	err = viper.Unmarshal(config)
	// Config file found and successfully parsed and returns a instance of Config type
	return config, nil
}
