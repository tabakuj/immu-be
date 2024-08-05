package main

import (
	"github.com/sirupsen/logrus"
	"immudb/internal"
	"immudb/internal/configuration"
)

// @title           Immudb Sample
// @version         1.0
// @description     This is a sample server that connects with immudb vault stores and receives data from there.

// @contact.name   Julian Tabaku
// @contact.email  julian.tabaku@outlook.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
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

	_, err = internal.NewServer(appConfigurations)
	return err
}
