package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	errors2 "immudb/internal/errors"
	"immudb/internal/services"
	"strconv"
)

type Handler struct {
	Service services.Service
	Engine  *gin.Engine
}

type Response struct {
	Data interface{} `json:"data"`
	Err  string      `json:"error_message"`
}

func NewHandler(bookService services.Service, router *gin.Engine) *Handler {
	handler := &Handler{
		Service: bookService,
		Engine:  router,
	}
	SetupHealth(router)
	v1 := router.Group("/v1/api")

	// register authors
	v1.GET("/account-info", handler.GetAccountInfos)
	v1.GET("/account-info/:id", handler.GetAccountInfo)
	v1.PUT("/account-info/:id", handler.UpdateAccountInfo)
	v1.POST("/account-info", handler.CreateAccountInfo)
	v1.DELETE("/account-info/:id", handler.DeleteAccountInfo)

	return handler
}

func AbortWithMessage(c *gin.Context, status int, err error, message string) {
	logrus.WithError(err).Error(message)

	// if custom validation error update status and message
	var badRequest *errors2.ServiceError
	if errors.As(err, &badRequest) {
		status = badRequest.Status
		message = err.Error()
	}

	c.AbortWithStatusJSON(status, Response{
		Err: message,
	})
}

func getParamUInt(c *gin.Context, paramName string) (uint, error) {
	id := c.Params.ByName(paramName)
	idValue, err := strconv.ParseUint(id, 10, 64)
	return uint(idValue), err
}

func getQueryParamUInt(c *gin.Context, paramName string) (int, error) {
	id := c.Query(paramName)
	idValue, err := strconv.ParseInt(id, 10, 64)
	return int(idValue), err
}

func returnOk(c *gin.Context, status int, data interface{}) {
	c.IndentedJSON(status, Response{
		Data: data,
	})
}
