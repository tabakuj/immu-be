package services

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	cerror "immudb/internal/errors"
	"immudb/internal/models"
	"immudb/internal/persistance"
	"net/http"
	"testing"
)

func TestAccountService_GetAllAccountInfos(t *testing.T) {
	tests := []struct {
		name            string
		accountDb       persistance.AccountDB
		isErrorExpected bool
		expectedError   error
		expectedResult  []*models.AccountInfo
	}{
		{
			name: "Test_Validity",
			accountDb: &MockAccountDB{
				GetAllCustomFunction: func(ctx context.Context, pageNr, pageSize int) ([]*models.AccountInfo, error) {
					return []*models.AccountInfo{
						{
							Id:      1,
							Name:    "as",
							Iban:    "a",
							Address: persistance.AddPointer("test"),
							Type:    persistance.AddPointer(models.Receiving),
							Amount:  0,
						},
					}, nil
				},
			},
			expectedResult: []*models.AccountInfo{
				{
					Id:      1,
					Name:    "as",
					Iban:    "a",
					Address: persistance.AddPointer("test"),
					Type:    persistance.AddPointer(models.Receiving),
					Amount:  0,
				},
			},
		},
		{
			name: "Test_Error_Handling",
			accountDb: &MockAccountDB{
				GetAllCustomFunction: func(ctx context.Context, pageNr, pageSize int) ([]*models.AccountInfo, error) {
					return nil, fmt.Errorf("test error")
				},
			},
			expectedError:   fmt.Errorf("test error"),
			isErrorExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			service, err := NewService(tt.accountDb)
			assert.NoError(t, err, "error setting up the service")

			// action
			infos, err := service.GetAllAccountInfos(context.Background(), 1, 10)

			//assert
			if tt.isErrorExpected {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				return
			}
			assert.NoError(t, err, "error loading account infos")
			assert.Equal(t, tt.expectedResult, infos)
		})
	}
}

func TestAccountService_GetAccountInfoById(t *testing.T) {
	tests := []struct {
		name            string
		accountDb       persistance.AccountDB
		requestId       uint
		isErrorExpected bool
		expectedError   error
		expectedResult  *models.AccountInfo
	}{
		{
			name: "Test_Validity",
			accountDb: &MockAccountDB{
				GetByCustomFunction: func(ctx context.Context, id uint) (*models.AccountInfo, error) {
					return &models.AccountInfo{
						Id:      1,
						Name:    "as",
						Iban:    "a",
						Address: persistance.AddPointer("test"),
						Type:    persistance.AddPointer(models.Receiving),
						Amount:  0,
					}, nil
				},
			},
			expectedResult: &models.AccountInfo{
				Id:      1,
				Name:    "as",
				Iban:    "a",
				Address: persistance.AddPointer("test"),
				Type:    persistance.AddPointer(models.Receiving),
				Amount:  0,
			},
		},
		{
			name: "Test_Error_Handling",
			accountDb: &MockAccountDB{
				GetByCustomFunction: func(ctx context.Context, id uint) (*models.AccountInfo, error) {
					return nil, fmt.Errorf("test error")
				},
			},
			expectedError:   fmt.Errorf("test error"),
			isErrorExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			service, err := NewService(tt.accountDb)
			assert.NoError(t, err, "error setting up the service")

			// action
			info, err := service.GetAccountInfoById(context.Background(), tt.requestId)

			//assert
			if tt.isErrorExpected {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				return
			}
			assert.NoError(t, err, "error loading account")
			assert.Equal(t, tt.expectedResult, info)
		})
	}
}

func TestAccountService_CreateAccountInfo(t *testing.T) {
	tests := []struct {
		name            string
		accountDb       persistance.AccountDB
		request         *models.AccountInfo
		isErrorExpected bool
		expectedError   error
		expectedResult  *models.AccountInfo
	}{
		{
			name: "Test_Validity",
			request: &models.AccountInfo{
				Name:    "as",
				Iban:    "a",
				Address: persistance.AddPointer("test"),
				Type:    persistance.AddPointer(models.Receiving),
				Amount:  11.1,
			},
			accountDb: &MockAccountDB{
				CreateCustomFunction: func(ctx context.Context, info models.AccountInfo) (*models.AccountInfo, error) {
					return &models.AccountInfo{
						Id:      1,
						Name:    "as",
						Iban:    "a",
						Address: persistance.AddPointer("test"),
						Type:    persistance.AddPointer(models.Receiving),
						Amount:  11.1,
					}, nil
				},
			},
			expectedResult: &models.AccountInfo{
				Id:      1,
				Name:    "as",
				Iban:    "a",
				Address: persistance.AddPointer("test"),
				Type:    persistance.AddPointer(models.Receiving),
				Amount:  11.1,
			},
		},
		{
			name: "Test_Error_Handling",
			request: &models.AccountInfo{
				Name:    "as",
				Iban:    "a",
				Address: persistance.AddPointer("test"),
				Type:    persistance.AddPointer(models.Receiving),
				Amount:  11.1,
			},
			accountDb: &MockAccountDB{
				CreateCustomFunction: func(ctx context.Context, info models.AccountInfo) (*models.AccountInfo, error) {
					return nil, fmt.Errorf("test error")
				},
			},
			expectedResult: &models.AccountInfo{
				Id:      1,
				Name:    "as",
				Iban:    "a",
				Address: persistance.AddPointer("test"),
				Type:    persistance.AddPointer(models.Receiving),
				Amount:  11.1,
			},
			expectedError:   fmt.Errorf("test error"),
			isErrorExpected: true,
		},
		{
			name:            "Test_Model_Validation_EmptyModel",
			request:         nil,
			accountDb:       &MockAccountDB{},
			expectedResult:  nil,
			expectedError:   cerror.NewServiceError("invalid input", http.StatusBadRequest),
			isErrorExpected: true,
		},
		{
			name: "Test_Model_Validation_IdSpecified",
			request: &models.AccountInfo{
				Id:      1,
				Name:    "as",
				Iban:    "a",
				Address: persistance.AddPointer("test"),
				Type:    persistance.AddPointer(models.Receiving),
				Amount:  11.1,
			},
			accountDb:       &MockAccountDB{},
			expectedResult:  nil,
			expectedError:   cerror.NewServiceError("invalid accountNumber for the account, do not specify accountNumber", http.StatusBadRequest),
			isErrorExpected: true,
		},
		{
			name: "Test_Model_Validation_NoName",
			request: &models.AccountInfo{
				Name:    "",
				Iban:    "a",
				Address: persistance.AddPointer("test"),
				Type:    persistance.AddPointer(models.Receiving),
				Amount:  11.1,
			},
			accountDb:       &MockAccountDB{},
			expectedResult:  nil,
			expectedError:   cerror.NewServiceError("invalid accountName for the account", http.StatusBadRequest),
			isErrorExpected: true,
		},
		{
			name: "Test_Model_Validation_NoIban",
			request: &models.AccountInfo{
				Name:    "asdf",
				Iban:    "",
				Address: persistance.AddPointer("test"),
				Type:    persistance.AddPointer(models.Receiving),
				Amount:  11.1,
			},
			accountDb:       &MockAccountDB{},
			expectedResult:  nil,
			expectedError:   cerror.NewServiceError("invalid iban for the account", http.StatusBadRequest),
			isErrorExpected: true,
		},
		{
			name: "Test_Model_Validation_NoType",
			request: &models.AccountInfo{
				Name:    "asdf",
				Iban:    "as",
				Address: persistance.AddPointer("test"),
				Amount:  11.1,
			},
			accountDb:       &MockAccountDB{},
			expectedResult:  nil,
			expectedError:   cerror.NewServiceError("invalid type for the account", http.StatusBadRequest),
			isErrorExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			service, err := NewService(tt.accountDb)
			assert.NoError(t, err, "error setting up the service")

			// action
			info, err := service.CreateAccountInfo(context.Background(), tt.request)

			//assert
			if tt.isErrorExpected {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				return
			}
			assert.NoError(t, err, "error creating account")
			assert.Equal(t, tt.expectedResult, info)
		})
	}
}
