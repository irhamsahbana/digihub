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

func (r *promotionRepository) UpdatePromotion(ctx context.Context, req *entity.UpdatePromotionRequest) error {
	var (
		args       = make([]any, 0)
		argsCount  = 0
		argsFilled = 0
	)

	if req.Title.Present && req.Title.Valid {
		argsCount++
	}
	if req.Link.Present && req.Link.Valid {
		argsCount++
	}
	if req.Image.Present && req.Image.Valid {
		argsCount++
	}

	query := `
		UPDATE promotions
		SET
		`

	if req.Title.Present && req.Title.Valid {
		query += " title = ?"
		args = append(args, req.Title.Val)
		argsFilled++
		if argsFilled < argsCount {
			query += ","
		}
	}

	if req.Link.Valid {
		query += " link = ?"
		args = append(args, req.Link.Val)
		argsFilled++
		if argsFilled < argsCount {
			query += ","
		}
	}

	if req.Image.Present && req.Image.Valid {
		query += " path = ?"
		args = append(args, req.Path)
		argsFilled++
		if argsFilled < argsCount {
			query += ","
		}
	}

	query += `
		WHERE
		id = ?
	`
	args = append(args, req.Id)

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), args...)
	if err != nil {
		req.RemoveImage()
		log.Error().Any("payload", req).Err(err).Msg("repo::UpdatePromotion - failed to update promotion")
		return err
	}

	return nil
}

func (r *promotionRepository) GetPromotionById(ctx context.Context, id string) (entity.Promotion, error) {
	data := entity.Promotion{}

	query := `
		SELECT id, title, path, link
		FROM promotions
		WHERE id = ?
	`

	err := r.db.GetContext(ctx, &data, r.db.Rebind(query), id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Str("id", id).Msg("repo::GetPromotionByID - promotion not found")
			return data, errmsg.NewCustomErrors(404).SetMessage("Promosi tidak ditemukan")
		}
		log.Error().Err(err).Str("id", id).Msg("repo::GetPromotionByID - failed to get promotion")
		return data, err
	}

	data.Image = config.Envs.App.BaseURL + "/" + strings.ReplaceAll(data.Path, "storage/", "api/storage/")
	return data, nil
}
