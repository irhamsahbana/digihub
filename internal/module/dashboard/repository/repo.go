package repository

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/dashboard/entity"
	"codebase-app/internal/module/dashboard/ports"
	"context"
	"strings"
	"unicode"

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

	type daoArea struct {
		Id    string `db:"id"`
		Area  string `db:"area"`
		Types string `db:"type"`
		Leads int    `db:"leads"`
		Key   string
	}

	var (
		res       entity.WACSummaryResponse
		potencies = make([]daoPotency, 0, 4)
		areas     = make([]daoArea, 0, 20)
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
		WHERE
			name != 'Used-car'
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

	// for used-car, the wo/do is based on wac that status is completed (invoice_number column is not null)
	query = `
		WITH total_wo_do_alt AS (
			SELECT
				COUNT(wac.id) AS total_wo_do
			FROM
				walk_around_checks wac
			WHERE
				wac.status = 'completed'
				AND Wac.is_used_car = TRUE
				AND wac.user_id = ?
				AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM') = ?
		),
		total_leads_alt AS (
			SELECT
				COUNT(wacc.id) AS total_leads
			FROM
				walk_around_check_conditions wacc
			LEFT JOIN
				walk_around_checks wac
				ON wac.id = wacc.walk_around_check_id
			WHERE
				wac.is_used_car = TRUE
				AND wac.user_id = ?
				AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM') = ?
		)
		SELECT
			'Used-car' AS title,
			(SELECT total_leads FROM total_leads_alt) AS total_potencial_leads,
			(SELECT total_leads FROM total_leads_alt) AS total_leads,
			(SELECT total_wo_do FROM total_wo_do_alt) AS total_wo_do
	`

	var summary entity.Summary
	err = r.db.QueryRowxContext(ctx, r.db.Rebind(query),
		req.UserId, req.Month,
		req.UserId, req.Month,
	).StructScan(&summary)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummary - failed to get wac summary")
		return res, err
	}

	res.Summaries = append(res.Summaries, summary)

	// make percentage from total leads
	/*
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
	*/

	var totalPercentage float64
	var maxIndex int

	// Calculate the distribution percentages and find the index of the maximum percentage
	for i, summary := range res.Summaries {
		// Initialize the percentage to zero
		var percentage float64

		// Calculate the percentage if the total distribution is not zero
		if res.TotalLeadDistributions != 0 {
			percentage = float64(summary.TotalLeads) / float64(res.TotalLeadDistributions) * 100
			percentage = float64(int(percentage*100)) / 100 // Round to two decimal places
		}

		// Append the percentage to the result
		res.DistributionOfLeads = append(res.DistributionOfLeads, entity.Distribution{
			Title:      summary.Title,
			Percentage: percentage,
		})

		// Track the total percentage sum
		totalPercentage += percentage

		// Track the index of the maximum percentage
		if res.DistributionOfLeads[maxIndex].Percentage < percentage {
			maxIndex = i
		}
	}

	// Adjust the total to ensure it sums to 100%
	difference := 100 - totalPercentage
	if len(res.DistributionOfLeads) > 0 {
		res.DistributionOfLeads[maxIndex].Percentage += difference
	}

	query = `
		SELECT
			id,
			name AS area,
			type
		FROM
			areas
	`

	err = r.db.SelectContext(ctx, &areas, r.db.Rebind(query))
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummary - failed to get wac summary")
		return res, err
	}

	queryArea := `
		SELECT
	`
	lastIndex := len(areas) - 1

	for idx, area := range areas {
		// replace all non-alphanumeric characters with underscore
		key := strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsNumber(r) {
				return r
			}
			return '_'
		}, area.Area)
		// to lower case
		key = strings.ToLower(key)
		areas[idx].Key = key

		queryArea += `
			COALESCE(SUM(CASE WHEN wacc.area_id = '` + area.Id + `' AND wacc.is_interested = TRUE THEN 1 ELSE 0 END), 0) AS ` + key + `
		`
		if idx != lastIndex {
			queryArea += ", "
		}
	}

	queryArea += `
		FROM
			walk_around_check_conditions wacc
		LEFT JOIN
			walk_around_checks wac
			ON wac.id = wacc.walk_around_check_id
		WHERE
			wac.user_id = ?
			AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM') = ?
	`

	TrendArea := make(map[string]any)

	err = r.db.QueryRowxContext(ctx, r.db.Rebind(queryArea), req.UserId, req.Month).MapScan(TrendArea)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummary - failed to get wac summary")
		return res, err
	}

	for _, area := range areas {
		key := area.Key
		leads, ok := TrendArea[key]
		if !ok {
			leads = 0
		}

		res.ServiceTrends = append(res.ServiceTrends, entity.Trend{
			Types: area.Types,
			Area:  area.Area,
			Leads: leads,
		})
	}

	return res, nil
}

func (r *dashboardRepository) GetWACSummaryTechnician(ctx context.Context, req *entity.WACSummaryRequest) (entity.TechWACSummaryResponse, error) {
	type daoPotency struct {
		Id    string `db:"id"`
		Name  string `db:"name"`
		Total int    `db:"total"`
	}

	var (
		res entity.TechWACSummaryResponse
	)
	res.Month = req.Month

	query := `
		WITH total_leads AS (
			SELECT
				COALESCE(SUM(1), 0) AS total_leads
			FROM
				walk_around_check_conditions wacc
			LEFT JOIN
				walk_around_checks wac
				ON wac.id = wacc.walk_around_check_id
			WHERE
				wac.branch_id = (SELECT branch_id FROM users WHERE id = ?)
				AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM') = ?
		)
		SELECT
			COALESCE(SUM(CASE WHEN wac.is_needs_follow_up = TRUE THEN 1 ELSE 0 END), 0) AS total_wac_need_follow_up,
			COALESCE(SUM(CASE WHEN wac.is_followed_up = TRUE THEN 1 ELSE 0 END), 0) AS total_wac_followed_up,
			(SELECT total_leads FROM total_leads) AS total_leads
		FROM
			walk_around_checks wac
		WHERE
			wac.branch_id = (SELECT branch_id FROM users WHERE id = ?)
			AND wac.status = 'completed'
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query),
		req.UserId, req.Month,
		req.UserId,
	).StructScan(&res)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummaryTechnician - failed to get wac summary")
		return res, err
	}

	query = `
		SELECT
			id,
			name,
			(
				SELECT
					COALESCE(SUM(1), 0)
				FROM
					walk_around_check_conditions wacc
				LEFT JOIN
					walk_around_checks wac
					ON wac.id = wacc.walk_around_check_id
				WHERE
					wacc.potency_id = potencies.id
					AND wac.branch_id = (SELECT branch_id FROM users WHERE id = ?)
					AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM') = ?
			) AS total
		FROM
			potencies
	`

	potencies := make([]daoPotency, 0, 4)
	err = r.db.SelectContext(ctx, &potencies, r.db.Rebind(query),
		req.UserId, req.Month,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummaryTechnician - failed to get wac summary")
		return res, err
	}

	/*
		for _, potency := range potencies {
			var percentage float64

			if res.TotalLeads != 0 {
				percentage = float64(potency.Total) / float64(res.TotalLeads) * 100
				percentage = float64(int(percentage*100)) / 100
			}

			res.DistributionOfLeads = append(res.DistributionOfLeads, entity.Distribution{
				Title:      potency.Name,
				Percentage: percentage,
			})
		}
	*/

	var totalPercentage float64
	var maxIndex int

	// Calculate the distribution percentages and find the index of the maximum percentage
	for i, potency := range potencies {
		var percentage float64

		if res.TotalLeads != 0 {
			percentage = float64(potency.Total) / float64(res.TotalLeads) * 100
			percentage = float64(int(percentage*100)) / 100 // Round to two decimal places
		}

		// Add to the result
		res.DistributionOfLeads = append(res.DistributionOfLeads, entity.Distribution{
			Title:      potency.Name,
			Percentage: percentage,
		})

		// Keep track of the total percentage
		totalPercentage += percentage

		// Track the index of the maximum percentage
		if res.DistributionOfLeads[maxIndex].Percentage < percentage {
			maxIndex = i
		}
	}

	// Adjust the total to ensure it sums to 100%
	difference := 100 - totalPercentage
	if len(res.DistributionOfLeads) > 0 {
		res.DistributionOfLeads[maxIndex].Percentage += difference
	}

	return res, nil
}
