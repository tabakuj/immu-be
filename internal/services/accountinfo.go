package services

import (
	"context"
	"immudb/internal/models"
)

func (s *service) CreateAccountInfo(ctx context.Context, data *models.AccountInfo) (*models.AccountInfo, error) {
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

func (s *service) GetAllAccountInfos(ctx context.Context, page, pageSize int) ([]*models.AccountInfo, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.Db.GetAllAccountInfos(ctx, page, pageSize)
}

func (s *service) GetAccountInfoById(ctx context.Context, Id uint) (*models.AccountInfo, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.Db.GetAccountInfoById(ctx, Id)
}

func (s *service) validateAccount(data *models.AccountInfo) error {
	if data == nil {
		return NewServiceError("invalid input")
	}

	if data.Id != 0 {
		return NewServiceError("invalid accountNumber for the account, do not specify accountNumber")
	}

	if data.Iban == "" {
		return NewServiceError("invalid iban for the account")
	}
	// we might need to add a logic to check if iban already exist into our system since its supposed to be unique
	// same logic might be considered for name

	if data.Name == "" {
		return NewServiceError("invalid accountName for the account")
	}

	if data.Type == nil {
		return NewServiceError("invalid type for the account")
	}
	return nil
}
