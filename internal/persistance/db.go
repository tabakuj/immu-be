package persistance

import (
	"context"
	"fmt"
	"immudb/internal/models"
)

type AccountDB interface {
	CreateAccountInfo(ctx context.Context, ata models.AccountInfo) (*models.AccountInfo, error)
	GetAllAccountInfos(ctx context.Context) ([]*models.AccountInfo, error)
	GetAccountInfoById(ctx context.Context, Id uint) (*models.AccountInfo, error)
}

type AccountImmmuDB struct {
	store map[uint]models.AccountInfo
}

func NewAccountDb() (*AccountImmmuDB, error) {
	repo := AccountImmmuDB{
		store: make(map[uint]models.AccountInfo),
	}
	return &repo, nil
}

func (db *AccountImmmuDB) GetAccountInfoById(ctx context.Context, Id uint) (*models.AccountInfo, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// check if we have any cancellation before continuing
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	raw, exist := db.store[Id]
	if !exist {
		return nil, fmt.Errorf("account id %s not found", Id)
	}

	return &raw, nil
}

func (db *AccountImmmuDB) GetAllAccountInfos(ctx context.Context) ([]*models.AccountInfo, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// check if we have any cancellation before continuing
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	var output []*models.AccountInfo
	for _, value := range db.store {
		output = append(output, &value)
	}

	return output, nil
}

func (db *AccountImmmuDB) CreateAccountInfo(ctx context.Context, ata models.AccountInfo) (*models.AccountInfo, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// check if we have any cancellation before continuing
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	ata.AccountNumber = uint(len(db.store) + 1)

	db.store[ata.AccountNumber] = ata
	return &ata, nil
}
