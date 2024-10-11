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
			u.name AS service_advisor_name,
			waca.status,
			waca.total_potential_leads,
			waca.total_leads,
			waca.total_revenue,
			waca.created_at
		FROM
			wac_activities waca
		LEFT JOIN users u ON waca.user_id = u.id
		WHERE
			1 = 1
	`

	if req.Search != "" {
		query += ` AND (waca.status ILIKE ? OR u.name ILIKE ?)`
		args = append(args, "%"+req.Search+"%", "%"+req.Search+"%")
	}

	if req.Date != "" {
		query += ` AND TO_CHAR(waca.created_at AT TIME ZONE '` + req.Timezone + `', 'YYYY-MM-DD') = ?`
		args = append(args, req.Date)
	}

	query += ` ORDER BY waca.created_at DESC`
	query += ` LIMIT ? OFFSET ?`
	args = append(args, req.Paginate, (req.Page-1)*req.Paginate)

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
