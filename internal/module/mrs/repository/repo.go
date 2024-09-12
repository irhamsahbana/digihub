package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/mrs/ports"

	"github.com/jmoiron/sqlx"
)

var _ ports.MRSRepository = &mrsRepository{}

type mrsRepository struct {
	db *sqlx.DB
}

func NewMRSRepository() *mrsRepository {
	return &mrsRepository{
		db: adapter.Adapters.DigihubPostgres,
	}
}
