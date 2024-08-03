package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"immudb/internal/models"
	"net/http"
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
	// in real life applications we need to consider using pagination
	//type QueryParameter struct {
	//	Limit  string `form:"limit,default=5" binding:"numeric"`
	//	Offset string `form:"offset,default=0" binding:"numeric"`
	//}

	result, err := h.Service.GetAllAccountInfos(c.Request.Context())
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
	id := c.Params.ByName("id")
	if id == "" {
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

func convertAccountInfoToDTO(sc *models.AccountInfo) *AccountInfoDto {
	if sc == nil {
		return nil
	}
	return &AccountInfoDto{
		//AccountNumber: "" + sc.AccountNumber,
		AccountName: sc.AccountName,
		Iban:        sc.Iban,
		Address:     sc.Address,
		Amount:      sc.Amount,
		Type:        sc.Type,
	}
}

func bindToAccountInfo(c *gin.Context) (*models.AccountInfo, error) {
	var input AccountInfoDto
	err := c.BindJSON(&input)
	if err != nil {
		return nil, err
	}
	authorData := models.AccountInfo{
		//AccountNumber: input.AccountNumber,
		AccountName: input.AccountName,
		Iban:        input.Iban,
		Address:     input.Address,
		Amount:      input.Amount,
		Type:        input.Type,
	}
	return &authorData, nil
}
