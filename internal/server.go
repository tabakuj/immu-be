package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"immudb/internal/handlers"
	"immudb/internal/persistance"
	"immudb/internal/services"
)

type Server struct {
	handler *handlers.Handler
}

func NewServer(port string) (*Server, error) {
	router := gin.Default()
	router.Use(handlers.CORSMiddleware())

	db, err := persistance.NewAccountDb()
	if err != nil {
		logrus.WithError(err).Fatal("couldn't connect to db")
	}
	accountService, err := services.NewService(db)
	if err != nil {
		logrus.WithError(err).Fatal("couldn't setup server")
	}
	handler := handlers.NewHandler(accountService, router)

	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		logrus.WithError(err).Errorf("Setting up service failed.")
		return nil, err
	}
	logrus.Infof("Application is running on port:%s", port)
	return &Server{handler: handler}, nil
}
