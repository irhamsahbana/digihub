package repository

import (
	"codebase-app/internal/module/dashboard/entity"
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

func (r *dashboardRepository) GetWACLineChart(ctx context.Context, req *entity.GetWACLineChartRequest) ([]entity.GetWACLineChartResponse, error) {
	var (
		res = make([]entity.GetWACLineChartResponse, 0)
	)

	query := `
		SELECT
			TO_CHAR(waca.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM-DD') AS date,
			COALESCE(SUM(
				CASE
					WHEN
						waca.status = 'offered'
					THEN waca.total_potential_leads ELSE 0 END
			), 0) AS total_potential_leads,
			COALESCE(SUM(
				CASE
					WHEN
						waca.status = 'wip'
					THEN waca.total_leads ELSE 0 END
			), 0) AS total_leads,
			COALESCE(SUM(
				CASE
					WHEN
						waca.status = 'completed'
					THEN waca.total_completed_leads ELSE 0 END
			), 0) AS total_completed_leads
		FROM
			wac_activities waca
		WHERE
			waca.created_at AT TIME ZONE 'Asia/Makassar'
			BETWEEN (TO_TIMESTAMP(?, 'YYYY-MM-DD') AT TIME ZONE 'UTC')
			AND (TO_TIMESTAMP(?, 'YYYY-MM-DD') AT TIME ZONE 'UTC' + time '23:59:59.999999')
		GROUP BY
			TO_CHAR(waca.created_at AT TIME ZONE 'Asia/Makassar', 'YYYY-MM-DD')
		`

	err := r.db.SelectContext(ctx, &res, r.db.Rebind(query), req.From, req.To)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repo::GetWACLineChart - failed to get wac line chart")
		return nil, err
	}

	// count expected days
	expectedDays := 0
	fromDate, err := time.Parse("2006-01-02", req.From)
	if err != nil {
		log.Error().Err(err).Any("from", req.From).Msg("failed to parse from date")
		return nil, err
	}

	toDate, err := time.Parse("2006-01-02", req.To)
	if err != nil {
		log.Error().Err(err).Any("to", req.To).Msg("failed to parse to date")
		return nil, err
	}

	listOfDates := make(map[string]bool)
	for fromDate.Before(toDate.AddDate(0, 0, 1)) { // add 1 day to toDate to include the last date
		expectedDays++
		listOfDates[fromDate.Format("2006-01-02")] = false
		fromDate = fromDate.AddDate(0, 0, 1)
	}

	// fill the data if blank for the expected days
	length := len(res)
	if length < expectedDays {
		newRes := make([]entity.GetWACLineChartResponse, 0)
		loc, _ := time.LoadLocation("Asia/Makassar")

		for i := 0; i < expectedDays; i++ {
			date := toDate.In(loc).AddDate(0, 0, -i).Format("2006-01-02")
			found := false
			for j := 0; j < length; j++ {
				if res[j].Date == date {
					newRes = append(newRes, res[j])
					listOfDates[date] = true
					found = true
					break
				}
			}
			if !found {
				newRes = append(newRes, entity.GetWACLineChartResponse{Date: date})
			}
		}

		res = newRes
	}

	// reverse the data
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}

	return res, nil
}
