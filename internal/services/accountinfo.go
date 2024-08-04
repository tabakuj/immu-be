package services

import (
	"context"
	"immudb/internal/errors"
	"immudb/internal/models"
	"net/http"
)

func (s *AccountService) CreateAccountInfo(ctx context.Context, data *models.AccountInfo) (*models.AccountInfo, error) {
	// cancellation check
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	err := s.validateAccount(data)
	if err != nil {
		return nil, err
	}
	return s.Db.CreateAccountInfo(ctx, *data)
}

func (s *AccountService) GetAllAccountInfos(ctx context.Context, page, pageSize int) ([]*models.AccountInfo, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.Db.GetAllAccountInfos(ctx, page, pageSize)
}

func (s *AccountService) GetAccountInfoById(ctx context.Context, Id uint) (*models.AccountInfo, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.Db.GetAccountInfoById(ctx, Id)
}

func (s *AccountService) validateAccount(data *models.AccountInfo) error {
	if data == nil {
		return errors.NewServiceError("invalid input", http.StatusBadRequest)
	}

	if data.Id != 0 {
		return errors.NewServiceError("invalid accountNumber for the account, do not specify accountNumber", http.StatusBadRequest)
	}

	if data.Iban == "" {
		return errors.NewServiceError("invalid iban for the account", http.StatusBadRequest)
	}
	// we might need to add a logic to check if iban already exist into our system since its supposed to be unique
	// same logic might be considered for name

	if data.Name == "" {
		return errors.NewServiceError("invalid accountName for the account", http.StatusBadRequest)
	}

	if data.Type == nil {
		return errors.NewServiceError("invalid type for the account", http.StatusBadRequest)
	}
	return nil
}
