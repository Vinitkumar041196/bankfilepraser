package config

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	GlobalConf          *viper.Viper `mapstructure:"-"`
	AppMode             string       `mapstructure:"APP_MODE"`
	PayRefRegex         string       `mapstructure:"PAYMENT_REFERENCE_REGEX"`
	FileColumnSeparator string       `mapstructure:"FILE_COLUMN_SEPARATOR"`
	DecimalPrecision    int          `mapstructure:"DECIMAL_PRECISION"`
	ServerAddress       string       `mapstructure:"SERVER_ADDR"`
}

func LoadConfig() (*AppConfig, error) {
	conf := viper.New()
	conf.AutomaticEnv()

	//load the config with env variables
	appConfig := new(AppConfig)
	appConfig.GlobalConf = conf
	appConfig.AppMode = conf.GetString("APP_MODE")
	appConfig.PayRefRegex = conf.GetString("PAYMENT_REFERENCE_REGEX")
	appConfig.FileColumnSeparator = conf.GetString("FILE_COLUMN_SEPARATOR")
	appConfig.DecimalPrecision = conf.GetInt("DECIMAL_PRECISION")
	appConfig.ServerAddress = conf.GetString("SERVER_ADDR")

	if appConfig.AppMode == "" { //if env variables not set check for .env file
		conf.SetConfigFile(".env")
		if err := conf.ReadInConfig(); err != nil {
			log.Printf("error reading config file, %v", err)
			return nil, err
		}
		err := conf.Unmarshal(&appConfig)
		if err != nil {
			log.Printf("error reading config file, %v", err)
			return nil, err
		}
	}

	return appConfig, nil
}
