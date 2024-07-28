package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/z_template_v2/ports"

	"github.com/jmoiron/sqlx"
)

type xxxRepository struct {
	db *sqlx.DB
}

func NewXxxRepository() ports.XxxRepository {
	return &xxxRepository{
		db: adapter.Adapters.DigihubPostgres,
	}
}
