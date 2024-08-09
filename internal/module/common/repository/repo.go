package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/common/entity"
	"codebase-app/internal/module/common/ports"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.CommonRepository = &commonRepository{}

type commonRepository struct {
	db *sqlx.DB
}

func NewCommonRepository() *commonRepository {
	return &commonRepository{
		db: adapter.Adapters.DigihubPostgres,
	}
}

func (r *commonRepository) GetAreas(ctx context.Context) ([]entity.CommonResponse, error) {
	var (
		result = make([]entity.CommonResponse, 0)
	)

	query := `
		SELECT
			id, name
		FROM
			areas
	`

	err := r.db.SelectContext(ctx, &result, query)
	if err != nil {
		log.Error().Err(err).Msg("repo::GetAreas - Failed to get areas")
		return nil, err
	}

	return result, nil
}

func (r *commonRepository) GetPotencies(ctx context.Context) ([]entity.CommonResponse, error) {
	var (
		result = make([]entity.CommonResponse, 0)
	)

	query := `
		SELECT
			id, name
		FROM
			potencies
	`

	err := r.db.SelectContext(ctx, &result, query)
	if err != nil {
		log.Error().Err(err).Msg("repo::GetPotencies - Failed to get potencies")
		return nil, err
	}

	return result, nil
}

func (r *commonRepository) GetVehicleTypes(ctx context.Context) ([]entity.CommonResponse, error) {
	var (
		result = make([]entity.CommonResponse, 0)
	)

	query := `
		SELECT
			id, name
		FROM
			vehicle_types
	`

	err := r.db.SelectContext(ctx, &result, query)
	if err != nil {
		log.Error().Err(err).Msg("repo::GetVehicleTypes - Failed to get vehicle types")
		return nil, err
	}

	return result, nil
}
