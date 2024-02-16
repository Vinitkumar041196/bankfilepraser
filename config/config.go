package config

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	PayRefRegex         string `mapstructure:"PAYMENT_REFERENCE_REGEX"`
	FilePath            string `mapstructure:"FILE_PATH"`
	FileColumnSeparator string `mapstructure:"FILE_COLUMN_SEPARATOR"`
	FileHasHeader       bool   `mapstructure:"FILE_HAS_HEADER"`
	DecimalPrecision    int    `mapstructure:"DECIMAL_PRECISION"`
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
