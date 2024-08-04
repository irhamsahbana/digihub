package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/wac/entity"
	"codebase-app/internal/module/wac/ports"
	"context"

	"github.com/jmoiron/sqlx"
)

var _ ports.WACRepository = &wacRepository{}

type wacRepository struct {
	db *sqlx.DB
}

func NewWACRepository() *wacRepository {
	return &wacRepository{
		db: adapter.Adapters.DigihubPostgres,
	}
}

func (r *wacRepository) CreateWAC(ctx context.Context, req *entity.CreateWACRequest) error {
	// your logic here

	return nil
}
