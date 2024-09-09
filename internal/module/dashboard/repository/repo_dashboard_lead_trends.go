package repository

import (
	"codebase-app/internal/module/dashboard/entity"
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

func (r *dashboardRepository) GetLeadsTrends(ctx context.Context, req *entity.LeadTrendsRequest) ([]entity.LeadTrendsResponse, error) {
	var res = make([]entity.LeadTrendsResponse, 0)

	query := `
		SELECT
			TO_CHAR(wacc.created_at, 'YYYY/Mon') AS month,
			SUM(
				CASE
					WHEN
						wac.user_id = ?
					THEN 1 ELSE 0 END
			) AS review_conditions,
			SUM(
				CASE
					WHEN
						wac.user_id = ?
						AND wacc.is_interested = TRUE
						AND wac.status != 'offered'
					THEN 1 ELSE 0 END
			) AS leads
		FROM
			walk_around_check_conditions wacc
		LEFT JOIN
			walk_around_checks wac
			ON wac.id = wacc.walk_around_check_id
		WHERE
			wacc.created_at >= NOW() - INTERVAL '11 months'
		GROUP BY
			TO_CHAR(wacc.created_at, 'YYYY/Mon')
		ORDER BY
			TO_CHAR(wacc.created_at, 'YYYY/Mon') DESC
	`

	err := r.db.SelectContext(ctx, &res, r.db.Rebind(query), req.UserId, req.UserId)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetLeadsTrends - failed to get leads trends")
		return nil, err
	}

	// feel the data if blank for 12 months

	length := len(res)

	if length < 12 {
		newRes := make([]entity.LeadTrendsResponse, 0)
		for i := 0; i < 12; i++ {
			month := time.Now().AddDate(0, -i, 0).Format("2006/Jan")
			found := false
			for j := 0; j < length; j++ {
				if res[j].Month == month {
					newRes = append(newRes, res[j])
					found = true
					break
				}
			}
			if !found {
				newRes = append(newRes, entity.LeadTrendsResponse{
					Month:            month,
					ReviewConditions: 0,
					Leads:            0,
				})
			}
		}

		res = newRes
	}

	// reverse the data
	// for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
	// 	res[i], res[j] = res[j], res[i]
	// }

	return res, nil
}
