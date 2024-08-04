package services

import (
	"context"
	"immudb/internal/models"
	"immudb/internal/persistance"
)

type Service interface {
	CreateAccountInfo(ctx context.Context, ata *models.AccountInfo) (*models.AccountInfo, error)
	GetAllAccountInfos(ctx context.Context, page, pageSize int) ([]*models.AccountInfo, error)
	GetAccountInfoById(ctx context.Context, Id uint) (*models.AccountInfo, error)
}

type AccountService struct {
	Db persistance.AccountDB
}

func NewService(db persistance.AccountDB) (*AccountService, error) {
	return &AccountService{Db: db}, nil
}
