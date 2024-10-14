package repository

import (
	"codebase-app/internal/module/dashboard/entity"
	"context"

	"github.com/rs/zerolog/log"
)

func (r *dashboardRepository) GetActivities(ctx context.Context, req *entity.GetActivitiesRequest) (entity.GetActivitiesResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.Activity
	}

	var (
		res  entity.GetActivitiesResponse
		data = make([]dao, 0)
		args = make([]any, 0)
	)
	res.Items = make([]entity.Activity, 0)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			waca.id,
			u.name AS employee_name,
			c.name AS client_name,
			b.name AS branch_name,
			c.vehicle_license_number,
			vt.name AS vehicle_type_name,
			c.phone,
			waca.status,
			waca.total_potential_leads,
			waca.total_leads,
			waca.total_revenue,
			waca.created_at
		FROM
			wac_activities waca
		LEFT JOIN
			users u
			ON waca.user_id = u.id
		LEFT JOIN
			walk_around_checks wac
			ON waca.wac_id = wac.id
		LEFT JOIN
			branches b
			ON wac.branch_id = b.id
		LEFT JOIN
			clients c
			ON wac.client_id = c.id
		LEFT JOIN
			vehicle_types vt
			ON c.vehicle_type_id = vt.id
		WHERE
			1 = 1
	`

	if req.Search != "" {
		query += ` AND (waca.status ILIKE ? OR u.name ILIKE ?)`
		args = append(args, "%"+req.Search+"%", "%"+req.Search+"%")
	}

	if req.From != "" && req.To != "" {
		query += ` AND waca.created_at AT TIME ZONE '` + req.Timezone + `'
		BETWEEN (TO_TIMESTAMP(?, 'YYYY-MM-DD') AT TIME ZONE 'UTC')
		AND (TO_TIMESTAMP(?, 'YYYY-MM-DD') AT TIME ZONE 'UTC' + time '23:59:59.999999')`
		args = append(args, req.From, req.To)
	}

	if req.BranchId != "" {
		query += ` AND wac.branch_id = ?`
		args = append(args, req.BranchId)

		queryBranchName := `SELECT name FROM branches WHERE id = ?`
		err := r.db.GetContext(ctx, &req.BranchName, r.db.Rebind(queryBranchName), req.BranchId)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repo::GetActivities - failed to get branch name")
			return res, err
		}
	}

	query += ` ORDER BY waca.created_at DESC`

	// convert int to bool
	isExport := req.Export == 1

	if !isExport {
		query += ` LIMIT ? OFFSET ?`
		args = append(args, req.Paginate, (req.Page-1)*req.Paginate)
	}

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), args...)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetActivities - failed to get activities")
		return res, err
	}

	if len(data) > 0 {
		res.Meta.TotalData = data[0].TotalData
	}

	for _, d := range data {
		res.Items = append(res.Items, d.Activity)
	}

	res.Meta.CountTotalPage(req.Page, req.Paginate, res.Meta.TotalData)
	return res, nil
}
