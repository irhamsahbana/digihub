package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/z_template_v2/ports"

	"github.com/jmoiron/sqlx"
)

var _ ports.XxxRepository = &xxxRepository{}

type xxxRepository struct {
	db *sqlx.DB
}

func NewXxxRepository() *xxxRepository {
	return &xxxRepository{
		db: adapter.Adapters.DigihubPostgres,
	}
}
