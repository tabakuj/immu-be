package configuration

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ApplicationConfiguration struct {
	DbConnectionString string
	Port               string
	ImmuDbUrl          string
	ImmudbSearchUrl    string
	ImmuDbApiKey       string // for simplicity i am taking this from here but normally with would be taken from env variable or some other  secret
}

func LoadConfiguration() (*ApplicationConfiguration, error) {
	v := viper.New()
	v.SetConfigName("default")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.Is(err, &configFileNotFoundError) {
			logrus.WithError(err).Warn("error loading config file")
		}
	}
	v.AutomaticEnv()
	var configurations ApplicationConfiguration
	err := v.UnmarshalExact(&configurations)
	if err != nil {
		return nil, err
	}
	return &configurations, nil
}
