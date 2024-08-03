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

func (s *service) GetAllAccountInfos(ctx context.Context) ([]*models.AccountInfo, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.Db.GetAllAccountInfos(ctx)
}

func (s *service) GetAccountInfoById(ctx context.Context, Id string) (*models.AccountInfo, error) {
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

	if data.AccountNumber != 0 {
		return NewServiceError("invalid accountNumber for the account")
	}

	if data.Iban != "" {
		return NewServiceError("invalid iban for the account")
	}

	if data.AccountName != "" {
		return NewServiceError("invalid accountName for the account")
	}

	if data.Type == nil {
		return NewServiceError("invalid type for the account")
	}

	return nil
}
