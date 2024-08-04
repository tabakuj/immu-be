package services

import (
	"context"
	"immudb/internal/models"
)

type MockAccountDB struct {
	CreateCustomFunction func(ctx context.Context, info models.AccountInfo) (*models.AccountInfo, error)
	GetAllCustomFunction func(ctx context.Context, pageNr, pageSize int) ([]*models.AccountInfo, error)
	GetByCustomFunction  func(ctx context.Context, id uint) (*models.AccountInfo, error)
}

func (m *MockAccountDB) CreateAccountInfo(ctx context.Context, info models.AccountInfo) (*models.AccountInfo, error) {
	return m.CreateCustomFunction(ctx, info)
}

func (m *MockAccountDB) GetAllAccountInfos(ctx context.Context, pageNr, pageSize int) ([]*models.AccountInfo, error) {
	return m.GetAllCustomFunction(ctx, pageNr, pageSize)
}

func (m *MockAccountDB) GetAccountInfoById(ctx context.Context, id uint) (*models.AccountInfo, error) {
	return m.GetByCustomFunction(ctx, id)
}
