package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/mrs/entity"
	"codebase-app/internal/module/mrs/ports"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
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

func (r *mrsRepository) GetMRSs(ctx context.Context, req *entity.GetMRSsRequest) (entity.GetMRSsResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.MRSItem
	}

	var (
		res  entity.GetMRSsResponse
		data = make([]dao, 0)
	)
	res.Items = make([]entity.MRSItem, 0)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			wac.id,
			c.name AS client,
			sa.name AS service_advisor,
			wac.follow_up_at
		FROM
			walk_around_checks wac
		LEFT JOIN
			clients c ON c.id = wac.client_id
		LEFT JOIN
			users sa ON sa.id = wac.user_id
		WHERE
			wac.branch_id = (SELECT branch_id FROM users WHERE id = ?)
			AND wac.is_needs_follow_up = TRUE
		ORDER BY
			wac.follow_up_at DESC
	`

	if err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), req.UserId); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetMRSs - Failed to get MRSs")
		return res, err
	}

	for _, d := range data {
		res.Items = append(res.Items, d.MRSItem)
	}

	if len(res.Items) > 0 {
		res.Meta.TotalData = data[0].TotalData
	}

	res.Meta.CountTotalPage(req.Page, req.Paginate, res.Meta.TotalData)

	return res, nil
}
