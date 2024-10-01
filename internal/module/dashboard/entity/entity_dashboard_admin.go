package entity

import (
	"codebase-app/pkg/errmsg"
	"time"

	"github.com/rs/zerolog/log"
)

/*
Admin WAC Line Chart - Start
*/
type GetWACLineChartRequest struct {
	From string `query:"from" validate:"datetime=2006-01-02"`
	To   string `query:"to" validate:"datetime=2006-01-02"`
}

func (r *GetWACLineChartRequest) Validate() error {
	fromDate, err := time.Parse("2006-01-02", r.From)
	if err != nil {
		log.Error().Err(err).Str("from", r.From).Msg("failed to parse from date")
		return err
	}

	toDate, err := time.Parse("2006-01-02", r.To)
	if err != nil {
		log.Error().Err(err).Str("to", r.To).Msg("failed to parse to date")
		return err
	}

	if fromDate.After(toDate) {
		log.Warn().Str("from", r.From).Str("to", r.To).Msg("from date is after to date")
		return errmsg.NewCustomErrors(400).
			Add("from", "tanggal awal tidak boleh lebih besar dari tanggal akhir").
			Add("to", "tanggal akhir tidak boleh lebih kecil dari tanggal awal")
	}

	return nil
}

type GetWACLineChartResponse struct {
	PotentialLeads int `json:"potential_leads"`
	Leads          int `json:"leads"`
	Completed      int `json:"completed"`
}

/*
	Admin WAC Line Chart - End
*/

/*
Admin Summary per month - Start

- Pie Diagrams (distribution of leads from service advisor and MRA)
- Distribution of leads based on area
*/
type GetSummaryPerMonthRequest struct {
	Month string `query:"month" validate:"datetime=2006-01"`
}

type GetSummaryPerMonthResponse struct {
	SASummary       []Summary      `json:"sa_summary"`  // based on area
	MRASummary      MRASummary     `json:"mra_summary"` //based on follow up
	SADistribution  []Distribution `json:"sa_distribution"`
	MRADistribution []Distribution `json:"mra_distribution"`
}

type MRASummary struct {
	TotalWACNeedFollowUp int `json:"total_wac_need_follow_up"`
	TotalWACFollowedUp   int `json:"total_wac_followed_up"`
}

/*
 Admin Summary per month - End
*/

/*
Admin SA Latest Activity - Start
*/
type GetSALatestActivityResponse struct {
	Items []SALatestActivity `json:"items"`
}

type SALatestActivity struct {
	Name                string `json:"name" db:"name"`
	Status              string `json:"status" db:"status"`
	TotalPotentialLeads int    `json:"total_potential_leads" db:"total_potential_leads"`
	TotalLeads          int    `json:"total_leads" db:"total_leads"`
	TotalRevenue        int    `json:"total_revenue" db:"total_revenue"`
}

/*
Admin SA Latest Activity - End
*/