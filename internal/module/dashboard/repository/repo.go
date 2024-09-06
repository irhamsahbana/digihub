package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/dashboard/entity"
	"codebase-app/internal/module/dashboard/ports"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
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

func (r *dashboardRepository) GetWACSummary(ctx context.Context, req *entity.WACSummaryRequest) (entity.WACSummaryResponse, error) {
	type daoPotency struct {
		Id   string `db:"id"`
		Name string `db:"name"`
	}

	var (
		res       entity.WACSummaryResponse
		potencies = make([]daoPotency, 0, 4)
	)
	res.Summaries = make([]entity.Summary, 0, 4)
	res.DistributionOfLeads = make([]entity.Distribution, 0, 4)
	res.Month = req.Month // 2006-01

	// the timestamp is timestamp with timezone, so we need to convert it to Asia/Makassar timezone
	query := `
		SELECT
			COUNT(wac.id) AS wac_counts
		FROM
			walk_around_checks wac
		WHERE
			wac.user_id = ?
			AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM') = ?
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.UserId, req.Month).Scan(&res.WACCounts)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummary - failed to get wac summary")
		return res, err
	}

	query = `
		SELECT
			id,
			name
		FROM
			potencies
	`
	err = r.db.SelectContext(ctx, &potencies, r.db.Rebind(query))
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummary - failed to get wac summary")
		return res, err
	}

	// total potensial leads -> kondisi offered sudah dihitung, kondisi yang interest dan tidak interest
	// total leads -> kondisi wip baru dihitung, kondisi yang interest
	// total wo/do -> kondisi completed baru dihitung, kondisi yang interest

	for _, potency := range potencies {
		query = `
			SELECT
				'` + potency.Name + `' AS title,
				COALESCE(SUM(1), 0) AS total_potencial_leads,
				COALESCE(SUM(CASE WHEN wacc.is_interested = TRUE AND wac.status != 'offered' THEN 1 ELSE 0 END), 0) AS total_leads,
				COALESCE(SUM(CASE WHEN wacc.is_interested = TRUE AND wac.status = 'completed' THEN 1 ELSE 0 END), 0) AS total_wo_do
			FROM
				walk_around_check_conditions wacc
			LEFT JOIN
				walk_around_checks wac
				ON wac.id = wacc.walk_around_check_id
			WHERE
				wac.user_id = ?
				AND wacc.potency_id = ?
				AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM') = ?
		`

		var summary entity.Summary
		err = r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.UserId, potency.Id, req.Month).StructScan(&summary)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummary - failed to get wac summary")
			return res, err
		}

		res.TotalLeadDistributions += summary.TotalLeads
		res.Summaries = append(res.Summaries, summary)
	}

	// make percentage from total leads
	for _, summary := range res.Summaries {
		// make percentage 2 decimal, make sure the total lead distribution is not 0
		var percentage float64

		if res.TotalLeadDistributions != 0 {
			percentage = float64(summary.TotalLeads) / float64(res.TotalLeadDistributions) * 100
			percentage = float64(int(percentage*100)) / 100
		}

		res.DistributionOfLeads = append(res.DistributionOfLeads, entity.Distribution{
			Title:      summary.Title,
			Percentage: percentage,
		})
	}

	return res, nil
}

func (r *dashboardRepository) GetWACSummaryTechnician(ctx context.Context, request *entity.WACSummaryRequest) (entity.WACSummaryResponse, error) {
	var (
		res entity.WACSummaryResponse
	)

	return res, nil
}
