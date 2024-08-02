package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"immudb/configuration"
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
	fmt.Println(appConfigurations)
	return nil
}
