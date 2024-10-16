package repository

import (
	"codebase-app/internal/infrastructure/config"
	"codebase-app/internal/module/dashboard/entity"
	"context"
	"strings"

	"github.com/rs/zerolog/log"
)

func (r *dashboardRepository) GetWACSummary(ctx context.Context, req *entity.WACSummaryRequest) (entity.WACSummaryResponse, error) {
	var (
		res entity.WACSummaryResponse
	)
	res.Summaries = make([]entity.Summary, 0, 4)
	res.DistributionOfLeads = make([]entity.Distribution, 0, 4)
	res.Month = req.Month // 2006-01

	// counting walk around checks based on user id and month
	err := r.summaryWACCount(ctx, req, &res)
	if err != nil {
		return res, err
	}

	// get walk around check summary per potency
	err = r.summaryPerPotency(ctx, req, &res)
	if err != nil {
		return res, err
	}

	err = r.summaryTiers(ctx, req, &res)
	if err != nil {
		return res, err
	}

	err = r.summaryPromotions(ctx, req, &res)
	if err != nil {
		return res, err
	}

	// get walk around check summary per area in percentage
	r.summaryLeadsDistribution(&res)

	// get walk around check summary per area (count leads area distribution)
	err = r.summaryWACArea(ctx, req, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *dashboardRepository) summaryWACCount(ctx context.Context, req *entity.WACSummaryRequest, res *entity.WACSummaryResponse) error {
	query := `
		SELECT
			COUNT(wac.id) AS wac_counts,
			COALESCE(SUM(CASE WHEN wac.status = 'offered' THEN 1 ELSE 0 END), 0)
				AS total_wac_on_offered
		FROM
			walk_around_checks wac
		WHERE
			wac.user_id = ?
			AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM') = ?
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.UserId, req.Month).Scan(&res.WACCounts, &res.TotalWACOnOffered)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummary - failed to get wac summary")
		return err
	}

	return nil
}

func (r *dashboardRepository) summaryPerPotency(ctx context.Context, req *entity.WACSummaryRequest, res *entity.WACSummaryResponse) error {
	type daoPotency struct {
		Id   string `db:"id"`
		Name string `db:"name"`
	}

	var potencies = make([]daoPotency, 0, 4)
	res.Summaries = make([]entity.Summary, 0, 4)

	query := `
		SELECT
			id,
			name
		FROM
			potencies
		WHERE
			name != 'Used-car'
	`
	err := r.db.SelectContext(ctx, &potencies, r.db.Rebind(query))
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummaryPerPotency - failed to get wac summary")
		return err
	}

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
			log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummaryPerPotency - failed to get wac summary")
			return err
		}

		res.TotalLeadDistributions += summary.TotalLeads
		res.Summaries = append(res.Summaries, summary)
	}

	// for used-car, the wo/do is based on w ac that status is completed (invoice_number column is not null)
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
				AND wac.status != 'offered'
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
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummaryPerPotency - failed to get wac summary")
		return err
	}

	res.Summaries = append(res.Summaries, summary)
	res.TotalLeadDistributions += summary.TotalLeads

	return nil
}

func (r *dashboardRepository) summaryLeadsDistribution(res *entity.WACSummaryResponse) {
	if res.TotalLeadDistributions == 0 {
		// Set all percentages to 0% since there are no leads
		for _, summary := range res.Summaries {
			res.DistributionOfLeads = append(res.DistributionOfLeads, entity.Distribution{
				Title:      summary.Title,
				Percentage: 0.0,
				Total:      summary.TotalLeads,
			})
		}
	} else {
		var totalPercentage float64
		var maxIndex int // Index of the highest percentage

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
				Total:      summary.TotalLeads,
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
	}
}

func (r *dashboardRepository) summaryWACArea(ctx context.Context, req *entity.WACSummaryRequest, res *entity.WACSummaryResponse) error {
	query := `
		SELECT
			a.name AS area,
			a.type,
			COALESCE(SUM(CASE WHEN wacc.area_id = a.id AND wacc.is_interested = TRUE THEN 1 ELSE 0 END), 0) AS leads
		FROM
			areas a
		LEFT JOIN
			walk_around_check_conditions wacc
			ON wacc.area_id = a.id
		LEFT JOIN
			walk_around_checks wac
			ON wac.id = wacc.walk_around_check_id
			AND wac.user_id = ?
			AND wac.status != 'offered'
			AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM') = ?
		GROUP BY
			a.type, a.name
		ORDER BY
			leads DESC, a.type, a.name ASC
		`

	err := r.db.SelectContext(ctx, &res.ServiceTrends, r.db.Rebind(query), req.UserId, req.Month)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummary - failed to get wac summary")
		return err
	}

	return nil
}

func (r *dashboardRepository) summaryTiers(ctx context.Context, req *entity.WACSummaryRequest, res *entity.WACSummaryResponse) error {
	var (
		currentYear            = res.Month[:4]
		totalRevenueNotUsedCar float64
		totalRevenueUsedCar    float64
		currentTier            string
		nextTier               *string
	)

	query := `
		SELECT
			COALESCE(SUM(CASE WHEN wac.status = 'completed' THEN wacc.revenue ELSE 0 END), 0) AS revenue
		FROM
			walk_around_check_conditions wacc
		LEFT JOIN
			walk_around_checks wac
			ON wac.id = wacc.walk_around_check_id
		WHERE
			wac.user_id = ?
			AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY') = ?
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.UserId, currentYear).Scan(&totalRevenueNotUsedCar)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummary - failed to get wac summary")
		return err
	}

	query = `
		SELECT
			COALESCE(
				SUM(
					CASE
						WHEN
							wac.status = 'completed'
							AND is_used_car = TRUE
						THEN wac.revenue
						ELSE 0
					END),
			0) AS revenue
		FROM
			walk_around_checks wac
		WHERE
			wac.user_id = ?
			AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY') = ?
	`

	err = r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.UserId, currentYear).Scan(&totalRevenueUsedCar)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummary - failed to get wac summary")
		return err
	}

	totalRevenue := totalRevenueNotUsedCar + totalRevenueUsedCar
	if totalRevenue >= 1500000 {
		currentTier = "platinum"
	} else if totalRevenue >= 1000000 {
		currentTier = "gold"
		nextTier = stringPointer("platinum")
	} else {
		currentTier = "silver"
		nextTier = stringPointer("gold")
	}

	res.Tiers = entity.Tier{
		Current: currentTier,
		Next:    nextTier,
		Revenue: totalRevenue,
	}

	return nil
}

func (r *dashboardRepository) summaryPromotions(ctx context.Context, req *entity.WACSummaryRequest, res *entity.WACSummaryResponse) error {
	type dao struct {
		Id   string `db:"id"`
		Path string `db:"path"`
	}

	promotions := make([]entity.Promotion, 0, 5)
	data := make([]dao, 0, 5)

	query := `
		SELECT
			p.id,
			p.path
		FROM
			promotions p
	`

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query))
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummary - failed to get wac summary")
		return err
	}

	for _, d := range data {
		promotions = append(promotions, entity.Promotion{
			Id:    d.Id,
			Image: config.Envs.App.BaseURL + "/" + strings.ReplaceAll(d.Path, "storage/", "api/storage/"),
		})
	}
	res.Promotions = promotions

	return nil
}

func stringPointer(s string) *string {
	return &s
}
