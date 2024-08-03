package main

import (
	"github.com/sirupsen/logrus"
	"immudb/internal"
	"immudb/internal/configuration"
)

func main() {

	err := setupApi()
	if err != nil {
		logrus.WithError(err).Error("failed to setup api")
	}
	logrus.Info("api setup complete")
}

func setupApi() error {
	appConfigurations, err := configuration.LoadConfiguration()
	if err != nil {
		return err
	}

	_, err = internal.NewServer(appConfigurations.Port)
	return err
}
