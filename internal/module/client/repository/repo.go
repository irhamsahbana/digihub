package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/client/entity"
	"codebase-app/internal/module/client/ports"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ClientRepository = &clientRepository{}

type clientRepository struct {
	db *sqlx.DB
}

func NewClientRepository() *clientRepository {
	return &clientRepository{
		db: adapter.Adapters.DigihubPostgres,
	}
}

func (r *clientRepository) GetClients(ctx context.Context, req *entity.GetClientsRequest) (entity.GetClientsResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.Client
	}

	var (
		data = make([]dao, 0)
		args = make([]any, 0)
		res  = entity.GetClientsResponse{}
	)
	res.Items = make([]entity.Client, 0)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			c.id,
			c.name,
			c.vehicle_license_number,
			vt.name AS vehicle_type,
			c.phone
		FROM clients c
		LEFT JOIN
			vehicle_types vt ON vt.id = c.vehicle_type_id
		WHERE
			1 = 1
	`

	if req.Search != "" {
		query += ` AND (c.name ILIKE ? OR c.vehicle_license_number ILIKE ? OR c.phone ILIKE ?) `
		args = append(args, "%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	query += ` ORDER BY c.name DESC LIMIT ? OFFSET ?`
	args = append(args, req.Paginate, (req.Page-1)*req.Paginate)

	if err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), args...); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetClients - failed to get clients")
		return res, err
	}

	for _, d := range data {
		res.Items = append(res.Items, d.Client)
	}

	if len(res.Items) > 0 {
		res.Meta.TotalData = data[0].TotalData
	}

	res.Meta.CountTotalPage(req.Page, req.Paginate, res.Meta.TotalData)
	return res, nil
}
