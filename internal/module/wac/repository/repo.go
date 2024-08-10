package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/wac/ports"

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
