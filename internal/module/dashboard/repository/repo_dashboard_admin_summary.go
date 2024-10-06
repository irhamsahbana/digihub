package repository

import (
	"codebase-app/internal/module/dashboard/entity"
	"context"

	"github.com/rs/zerolog/log"
)

func (r *dashboardRepository) GetAdminSummary(ctx context.Context, req *entity.GetSummaryPerMonthRequest) (entity.GetSummaryPerMonthResponse, error) {
	var res entity.GetSummaryPerMonthResponse
	res.SADistribution = make([]entity.Distribution, 0, 4)
	res.MRADistribution = make([]entity.Distribution, 0, 2)

	err := r.getSASummary(ctx, req, &res)
	if err != nil {
		return res, err
	}

	err = r.getMRASummary(ctx, req, &res)
	if err != nil {
		return res, err
	}

	// err = r.getSADistribution(ctx, req, &res)
	// if err != nil {
	// 	return res, err
	// }

	return res, nil
}

func (r *dashboardRepository) getSASummary(ctx context.Context, req *entity.GetSummaryPerMonthRequest, res *entity.GetSummaryPerMonthResponse) error {
	type daoPotency struct {
		Id   string `db:"id"`
		Name string `db:"name"`
	}

	var potencies = make([]daoPotency, 0, 4)
	res.SASummary = make([]entity.Summary, 0, 4)

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
		log.Error().Err(err).Any("payload", req).Msg("repo::getSASummary - failed to get wac summary")
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
				wacc.potency_id = ?
				AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM') = ?
		`

		var summary entity.Summary
		err = r.db.QueryRowxContext(ctx, r.db.Rebind(query), potency.Id, req.Month).StructScan(&summary)
		if err != nil {
			log.Error().Err(err).Any("payload", req).Msg("repo::getSASummary - failed to get wac summary")
			return err
		}

		res.SASummary = append(res.SASummary, summary)
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
				AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM') = ?
		)
		SELECT
			'Used-car' AS title,
			(SELECT total_leads FROM total_leads_alt) AS total_potencial_leads,
			(SELECT total_leads FROM total_leads_alt) AS total_leads,
			(SELECT total_wo_do FROM total_wo_do_alt) AS total_wo_do
	`

	var summary entity.Summary
	err = r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.Month, req.Month).StructScan(&summary)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::getSASummary - failed to get wac summary")
		return err
	}

	res.SASummary = append(res.SASummary, summary)

	return nil
}

func (r *dashboardRepository) getMRASummary(ctx context.Context, req *entity.GetSummaryPerMonthRequest, res *entity.GetSummaryPerMonthResponse) error {
	query := `
	WITH total_leads_wac_need_follow_up AS (
		SELECT
			COALESCE(SUM(1), 0) AS total_leads
		FROM
			walk_around_check_conditions wacc
		LEFT JOIN
			walk_around_checks wac
			ON wac.id = wacc.walk_around_check_id
		WHERE
			wacc.is_interested = TRUE
			AND wac.status = 'completed'
			AND wac.is_needs_follow_up = TRUE
			AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM') = ?
	)
	SELECT
		COALESCE(SUM(CASE WHEN wac.is_needs_follow_up = TRUE THEN 1 ELSE 0 END), 0) AS total_wac_need_follow_up,
		COALESCE(SUM(CASE WHEN wac.is_followed_up = TRUE THEN 1 ELSE 0 END), 0) AS total_wac_followed_up,
		(SELECT total_leads FROM total_leads_wac_need_follow_up) AS total_leads
	FROM
		walk_around_checks wac
	WHERE
		wac.status = 'completed'
		AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM') = ?
`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.Month, req.Month).StructScan(&res.MRASummary)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummaryTechnician - failed to get wac summary")
		return err
	}

	return nil
}
