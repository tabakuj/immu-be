package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"immudb/internal/models"
	"net/http"
)

const (
	DefaultPageSize = 100
	DefaultPage     = 1
)

type AccountInfoDto struct {
	AccountNumber uint                `json:"account_number"`
	AccountName   string              `json:"account_name"`
	Iban          string              `json:"iban"`
	Address       *string             `json:"address"`
	Amount        float64             `json:"amount"`
	Type          *models.AccountType `json:"type"`
}

func (h *Handler) GetAccountInfos(c *gin.Context) {
	page, err := getQueryParamUInt(c, "page")
	if err != nil {
		page = DefaultPage
	}
	pageSize, err := getQueryParamUInt(c, "pageSize")
	if err != nil {
		pageSize = DefaultPageSize
	}

	result, err := h.Service.GetAllAccountInfos(c.Request.Context(), page, pageSize)
	if err != nil {
		AbortWithMessage(c, http.StatusInternalServerError, err, "failed to load accountInfos")
		return
	}
	var output []*AccountInfoDto
	for _, item := range result {
		output = append(output, convertAccountInfoToDTO(item))
	}
	returnOk(c, http.StatusOK, output)
}

func (h *Handler) GetAccountInfo(c *gin.Context) {
	id, err := getParamUInt(c, "id")
	if err != nil {
		AbortWithMessage(c, http.StatusBadRequest, fmt.Errorf("please specify id"), "id is required")
		return
	}
	result, err := h.Service.GetAccountInfoById(c.Request.Context(), id)
	if err != nil {
		AbortWithMessage(c, http.StatusInternalServerError, err, "failed to load accountIndo")
		return
	}
	returnOk(c, http.StatusOK, convertAccountInfoToDTO(result))
}

func (h *Handler) DeleteAccountInfo(c *gin.Context) {
	AbortWithMessage(c, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"), "delete is not allowed")
}

func (h *Handler) CreateAccountInfo(c *gin.Context) {
	// prepare input
	input, err := bindToAccountInfo(c)
	if err != nil {
		AbortWithMessage(c, http.StatusBadRequest, err, "error binding to json")
		return
	}
	// execute
	data, err := h.Service.CreateAccountInfo(c.Request.Context(), input)
	if err != nil {
		AbortWithMessage(c, http.StatusInternalServerError, err, "failed to create author")
		return
	}
	returnOk(c, http.StatusCreated, convertAccountInfoToDTO(data))
}

func (h *Handler) UpdateAccountInfo(c *gin.Context) {
	AbortWithMessage(c, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"), "update  is not allowed")
}

func bindToAccountInfo(c *gin.Context) (*models.AccountInfo, error) {
	var input AccountInfoDto
	err := c.BindJSON(&input)
	if err != nil {
		return nil, err
	}
	return convertDtoToAccountInfo(&input), nil
}
