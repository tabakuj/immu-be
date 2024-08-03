package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"immudb/internal/configuration"
	"immudb/internal/handlers"
	"immudb/internal/persistance"
	"immudb/internal/services"
)

type Server struct {
	handler *handlers.Handler
}

func NewServer(config *configuration.ApplicationConfiguration) (*Server, error) {
	router := gin.Default()
	router.Use(handlers.CORSMiddleware())

	db := persistance.NewImmmuDB(config.ImmuDbUrl, config.ImmuDbApiKey, config.ImmudbSearchUrl, config.ImmuDbApiReadKey)

	accountService, err := services.NewService(db)
	if err != nil {
		logrus.WithError(err).Fatal("couldn't setup server")
	}
	handler := handlers.NewHandler(accountService, router)

	err = router.Run(fmt.Sprintf(":%s", config.Port))
	if err != nil {
		logrus.WithError(err).Errorf("Setting up service failed.")
		return nil, err
	}
	logrus.Infof("Application is running on port:%s", config.Port)
	return &Server{handler: handler}, nil
}
