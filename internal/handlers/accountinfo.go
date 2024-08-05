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

// GetAccountInfos lists all existing accounts
//
// @Summary      Get all account infos
// @Description  Retrieve a paginated list of all account information
// @ID           get-all-account-infos
// @Accept       json
// @Produce      json
// @Param        page     query    int    false  "Page number"       default(1)
// @Param        pageSize query    int    false  "Page size"         default(10)
// @Success      200      {object}   Response[[]AccountInfoDto] "List of account info"
// @Failure      500      {object} Response[string]  "Internal server error"
// @Router       /account-info [get]
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

// GetAccountInfo
//
// @Summary      Get account info by ID
// @Description  Retrieve account information for a specific account by ID
// @ID           get-account-info-by-id
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  Response[AccountInfoDto] "Account information"
// @Failure      400  {object}  Response[string]  "Bad request, ID is required"
// @Failure      500  {object}  Response[string]  "Internal server error"
// @Router       /account-info/{id} [get]
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

// CreateAccountInfo
//
// @Summary      Create account info
// @Description  Add a new account information entry
// @ID           create-account-info
// @Accept       json
// @Produce      json
// @Param        account  body      AccountInfoDto  true  "Account information"
// @Success      201      {object}   Response[AccountInfoDto] "Created account info"
// @Failure      400      {object}  Response[string]  "Bad request, invalid input"
// @Failure      500      {object}  Response[string]  "Internal server error"
// @Router       /account-info [post]
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

// DeleteAccountInfo
//
// @Summary      Delete account info
// @Description  Delete account information for a specific account by ID (currently not allowed)
// @ID           delete-account-info
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Account ID"
// @Failure      405 {object}  Response[string] "Method not allowed"
// @Router       /account-info/{id} [delete]
func (h *Handler) DeleteAccountInfo(c *gin.Context) {
	AbortWithMessage(c, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"), "delete is not allowed")
}

// UpdateAccountInfo
//
// @Summary      Update account info
// @Description  Update account information for a specific account by ID (currently not allowed)
// @ID           update-account-info
// @Accept       json
// @Produce      json
// @Param        id      path      int          true  "Account ID"
// @Param        account body      AccountInfoDto  true  "Updated account information"
// @Failure      405     {object}  Response[string] "Method not allowed"
// @Router       /account-info/{id} [put]
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
