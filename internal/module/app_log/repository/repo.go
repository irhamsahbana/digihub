package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/app_log/entity"
	"codebase-app/internal/module/app_log/ports"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.AppLogRepository = &appLogRepository{}

type appLogRepository struct {
	db *sqlx.DB
}

func NewAppLogRepository() *appLogRepository {
	return &appLogRepository{
		db: adapter.Adapters.DigihubPostgres,
	}
}

func (r *appLogRepository) CreateLog(ctx context.Context, req *entity.CreateLogRequest) error {
	query := `
		INSERT INTO app_logs (
			app_id,
			log_level,
			info,
			message
		) VALUES (?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), req.AppId, req.LogLevel, req.Info, req.Message)

	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::CreateLog - An error occurred")
		return err
	}

	return nil
}
