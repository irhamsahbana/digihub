package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/wac/entity"
	"codebase-app/internal/module/wac/ports"
	"context"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
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

func (r *wacRepository) GetWACs(ctx context.Context, req *entity.GetWACsRequest) (entity.GetWACsResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.WacItem
	}

	var (
		query = strings.Builder{}
		args  = make([]any, 0, 8)
		res   entity.GetWACsResponse
		data  = make([]dao, 0, req.Paginate)
	)
	res.Items = make(map[string][]entity.WacItem)

	query.WriteString(`
		SELECT
			COUNT(*) OVER() AS total_data,
			wac.id,
			c.name AS client_name,
			wac.status,
			wac.created_at
		FROM
			walk_around_checks wac
		LEFT JOIN
			walk_around_check_conditions wacc ON wacc.walk_around_check_id = wac.id
		LEFT JOIN
			clients c ON c.id = wac.client_id
		WHERE
			wac.deleted_at IS NULL
			AND (wac.user_id = ? OR wacc.assigned_user_id = ?)
	`)
	args = append(args, req.UserId, req.UserId)

	if req.Query != "" {
		query.WriteString(" AND (c.name ILIKE ? OR c.vehicle_license_number ILIKE ?)")
		args = append(args, "%"+req.Query+"%", "%"+req.Query+"%")
	}

	if req.Status != "" {
		query.WriteString(" AND wac.status = ?")
		args = append(args, req.Status)
	}

	query.WriteString(" ORDER BY wac.created_at DESC")

	query.WriteString(" LIMIT ? OFFSET ?")
	args = append(args, req.Paginate, (req.Page-1)*req.Paginate)

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query.String()), args...)
	if err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("repo::GetWACs - failed to get wacs")
		return res, err
	}

	// today := time.Now().In(time.FixedZone("WITA", 28800)).Format("2006-01-02") // WITA
	today := time.Now().Format("2006-01-02")
	length := len(data)
	for length > 0 {
		length--
		d := data[length]

		// date := d.CreatedAt.In(time.FixedZone("WITA", 28800)).Format("2006-01-02")
		date := d.CreatedAt.Format("2006-01-02")
		if today == date {
			date = "Hari ini"
		}

		res.Items[date] = append(res.Items[date], d.WacItem)
	}

	if len(data) > 0 {
		res.Meta.TotalData = data[0].TotalData
	}

	res.Meta.CountTotalPage(req.Page, req.Paginate, res.Meta.TotalData)

	return res, nil
}
