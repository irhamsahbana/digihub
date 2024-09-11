package repository

import (
	"codebase-app/internal/module/dashboard/entity"
	"context"

	"github.com/rs/zerolog/log"
)

func (r *dashboardRepository) GetWACSummaryTechnician(ctx context.Context, req *entity.WACSummaryRequest) (entity.TechWACSummaryResponse, error) {
	var (
		res entity.TechWACSummaryResponse
	)
	res.Month = req.Month

	// get walk around check summary that needs follow up
	err := r.summaryTechnicianNeedFollowUp(ctx, req, &res)
	if err != nil {
		return res, err
	}

	// get walk around check summary per potency
	potencies, err := r.summaryTechnicianTotalLeadsPerPotency(ctx, req)
	if err != nil {
		return res, err
	}

	// get walk around check summary per potency in percentage
	r.summaryTechnicianPercentagesOnLeadsDistributionPerPotency(potencies, &res)

	return res, nil
}

// this function will get the total leads from walk around check conditions that needs follow up
// also counting walk around checks that needs follow up and already followed up
// based on the user branch, month in Asia/Makassar timezone

func (r *dashboardRepository) summaryTechnicianNeedFollowUp(ctx context.Context, req *entity.WACSummaryRequest, res *entity.TechWACSummaryResponse) error {
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
				AND wac.branch_id = (SELECT branch_id FROM users WHERE id = ?)
				AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM') = ?
		)
		SELECT
			COALESCE(SUM(CASE WHEN wac.is_needs_follow_up = TRUE THEN 1 ELSE 0 END), 0) AS total_wac_need_follow_up,
			COALESCE(SUM(CASE WHEN wac.is_followed_up = TRUE THEN 1 ELSE 0 END), 0) AS total_wac_followed_up,
			(SELECT total_leads FROM total_leads_wac_need_follow_up) AS total_leads
		FROM
			walk_around_checks wac
		WHERE
			wac.branch_id = (SELECT branch_id FROM users WHERE id = ?)
			AND wac.status = 'completed'
			AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM') = ?
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query),
		req.UserId, req.Month,
		req.UserId, req.Month,
	).StructScan(res)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummaryTechnician - failed to get wac summary")
		return err
	}

	return nil
}

type daoTotalLeadsPerPotency struct {
	Id    string `db:"id"`
	Name  string `db:"name"`
	Total int    `db:"total"`
}

// this function will get the total leads per potency
// from walk around check conditions that needs follow up
// based on the user branch, month in Asia/Makassar timezone
// leads (condition that interested and walk around check status is completed)
func (r *dashboardRepository) summaryTechnicianTotalLeadsPerPotency(ctx context.Context, req *entity.WACSummaryRequest) ([]daoTotalLeadsPerPotency, error) {
	query := `
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
				AND wac.status = 'completed'
				AND wac.is_needs_follow_up = TRUE
				AND wacc.is_interested = TRUE
				AND wac.branch_id = (SELECT branch_id FROM users WHERE id = ?)
				AND TO_CHAR(wac.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM') = ?
		) AS total
	FROM
		potencies
`

	potencies := make([]daoTotalLeadsPerPotency, 0, 4)
	err := r.db.SelectContext(ctx, &potencies, r.db.Rebind(query),
		req.UserId, req.Month,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACSummaryTechnician - failed to get wac summary")
		return potencies, err
	}

	return potencies, nil
}

// this function will calculate the percentage of leads distribution per potency
func (r *dashboardRepository) summaryTechnicianPercentagesOnLeadsDistributionPerPotency(potencies []daoTotalLeadsPerPotency, res *entity.TechWACSummaryResponse) {
	if res.TotalLeads == 0 {
		// Set all percentages to 0% since there are no leads
		for _, potency := range potencies {
			res.DistributionOfLeads = append(res.DistributionOfLeads, entity.Distribution{
				Title:      potency.Name,
				Percentage: 0.0,
			})
		}
	} else {
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
	}
}
