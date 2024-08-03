package persistance

import (
	"context"
	"immudb/internal/models"
)

type AccountDB interface {
	CreateAccountInfo(ctx context.Context, ata models.AccountInfo) (*models.AccountInfo, error)
	GetAllAccountInfos(ctx context.Context, pageNr, pageSize int) ([]*models.AccountInfo, error)
	GetAccountInfoById(ctx context.Context, Id uint) (*models.AccountInfo, error)
}
