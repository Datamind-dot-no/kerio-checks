package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config represent the application configuration
type Config struct {
	ServerName         string `mapstructure:"server_name"`
	KerioStorePath     string `mapstructure:"keriostore_path"`
	Sender             string `mapstructure:"sender"`
	Recipient          string `mapstructure:"recipient"`
	SubjectT           string `mapstructure:"subject_t"`
	HTMLBodyT          string `mapstructure:"html_body"`
	TextBodyT          string `mapstructure:"text_body"`
	CharSet            string `mapstructure:"char_set"`
	QueuePath          string `mapstructure:"q_path"`
	QueueWarnThreshold int    `mapstructure:"q_warn_threshold"`
}

/*
// Sadly, it doesn't work when using json tags, must use mapstructure instead, and NOT have whitespace after - all hail https://github.com/spf13/viper/issues/498#issuecomment-410485531
type Config struct {
	ServerName         string `json:"server_name"`
	KerioStorePath     string `json:"keriostore_path"`
	Sender             string `json:"sender"`
	Recipient          string `json:"recipient"`
	SubjectT           string `json:"subject_t"`
	HTMLBodyT          string `json:"html_body"`
	TextBodyT          string `json:"text_body"`
	CharSet            string `json:"char_set"`
	QueuePath          string `json:"q_path"`
	QueueWarnThreshold int    `json:"q_warn_threshold"`
}
*/

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

	var C Config
	err = viper.Unmarshal(&C)
	if err != nil {
		fmt.Println("unable to decode into struct")
		return &C, err
	}

	// debug statements
	//fmt.Println("The Server configured is: " + C.ServerName)
	//fmt.Println("The Kerio Store path configured is: " + C.KerioStorePath)
	//fmt.Println("The queue warning threshold is now configured as: " + strconv.Itoa(C.QueueWarnThreshold))

	// Config file found and successfully parsed and returns a instance of Config type
	return &C, nil
}
