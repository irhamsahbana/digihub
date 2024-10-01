package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/infrastructure/config"
	"codebase-app/internal/module/promotion/entity"
	"codebase-app/internal/module/promotion/ports"
	"codebase-app/pkg/errmsg"
	"context"
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

var _ ports.PromotionRepository = &promotionRepository{}

type promotionRepository struct {
	db *sqlx.DB
}

func NewPromotionRepository() *promotionRepository {
	return &promotionRepository{
		db: adapter.Adapters.DigihubPostgres,
	}
}

func (r *promotionRepository) CreatePromotion(ctx context.Context, req *entity.CreatePromotionRequest) error {
	query := `
		INSERT INTO promotions (id, title, path, link)
		VALUES (?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query),
		ulid.Make().String(), req.Title, req.Path, req.Link,
	)
	if err != nil {
		req.RemoveImage()
		log.Error().Any("payload", req).Err(err).Msg("repo::CreatePromotion - failed to insert promotion")
		return err
	}

	return nil
}

func (r *promotionRepository) GetPromotions(ctx context.Context) ([]entity.Promotion, error) {
	data := make([]entity.Promotion, 0)

	query := `
		SELECT id, title, path, link
		FROM promotions
	`

	err := r.db.SelectContext(ctx, &data, query)
	if err != nil {
		log.Error().Err(err).Msg("repo::GetPromotions - failed to get promotions")
		return nil, err
	}

	for i := range data {
		data[i].Image = config.Envs.App.BaseURL + "/" + strings.ReplaceAll(data[i].Path, "storage/", "api/storage/")
	}

	return data, nil
}

func (r *promotionRepository) DeletePromotion(ctx context.Context, req *entity.DeletePromotionRequest) error {
	query := `
		DELETE FROM promotions
		WHERE id = ?
		RETURNING path
	`

	err := r.db.GetContext(ctx, &req.Path, r.db.Rebind(query), req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Any("payload", req).Msg("repo::DeletePromotion - promotion not found")
			return errmsg.NewCustomErrors(404).SetMessage("Promosi tidak ditemukan")
		}
		log.Error().Err(err).Any("payload", req).Msg("repo::DeletePromotion - failed to delete promotion")
		return err
	}

	req.Path = strings.ReplaceAll(req.Path, "storage/", "./storage/")

	return nil
}
