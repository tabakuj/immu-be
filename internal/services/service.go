package services

import (
	"context"
	"immudb/internal/models"
	"immudb/internal/persistance"
)

type Service interface {
	CreateAccountInfo(ctx context.Context, ata *models.AccountInfo) (*models.AccountInfo, error)
	GetAllAccountInfos(ctx context.Context) ([]*models.AccountInfo, error)
	GetAccountInfoById(ctx context.Context, Id uint) (*models.AccountInfo, error)
}

type service struct {
	Db persistance.AccountDB
}

func NewService(db persistance.AccountDB) (*service, error) {
	return &service{Db: db}, nil
}
