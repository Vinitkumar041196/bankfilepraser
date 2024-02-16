package config

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	AppMode             string `mapstructure:"APP_MODE"`
	PayRefRegex         string `mapstructure:"PAYMENT_REFERENCE_REGEX"`
	FileColumnSeparator string `mapstructure:"FILE_COLUMN_SEPARATOR"`
	FileHasHeader       bool   `mapstructure:"FILE_HAS_HEADER"`
	DecimalPrecision    int    `mapstructure:"DECIMAL_PRECISION"`
	ServerAddress       string `mapstructure:"SERVER_ADDR"`
	EnableTLS           bool   `mapstructure:"ENABLE_TLS"`
	SSLCertPath         string `mapstructure:"SSL_CRT_PATH"`
	SSLKeyPath          string `mapstructure:"SSL_KEY_PATH"`
}

func LoadConfig() (*AppConfig, error) {
	conf := viper.New()
	conf.SetConfigFile(".env")
	conf.AutomaticEnv()

	if err := conf.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %v", err)
		return nil, err
	}

	appConfig := new(AppConfig)
	err := conf.Unmarshal(&appConfig)
	if err != nil {
		log.Printf("Error reading config file, %v", err)
		return nil, err
	}

	return appConfig, nil
}
