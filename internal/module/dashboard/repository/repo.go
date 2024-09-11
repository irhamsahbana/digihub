package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/dashboard/ports"

	"github.com/jmoiron/sqlx"
)

var _ ports.DashboardRepository = &dashboardRepository{}

type dashboardRepository struct {
	db *sqlx.DB
}

func NewDashboardRepository() *dashboardRepository {
	return &dashboardRepository{
		db: adapter.Adapters.DigihubPostgres,
	}
}
