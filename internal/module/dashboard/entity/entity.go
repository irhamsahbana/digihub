package entity

import (
	"codebase-app/pkg/types"
	"time"
)

type LeadTrendsRequest struct {
	UserId string
}

type LeadTrendsResponse struct {
	Month            string `json:"month" db:"month"`
	ReviewConditions int    `json:"review_conditions" db:"review_conditions"`
	Leads            int    `json:"leads" db:"leads"`
}

type TechWACSummaryResponse struct {
	Month                string         `json:"month"`
	TotalWACNeedFollowUp int            `json:"total_wac_need_follow_up" db:"total_wac_need_follow_up"`
	TotalWACFollowedUp   int            `json:"total_wac_followed_up" db:"total_wac_followed_up"`
	TotalLeads           int            `json:"total_leads" db:"total_leads"`
	DistributionOfLeads  []Distribution `json:"distribution_of_leads"`
}

type GetActivitiesRequest struct {
	Page     int    `query:"page" validate:"required"`
	Paginate int    `query:"paginate" validate:"required"`
	Search   string `query:"search" validate:"omitempty,min=3"`
	Date     string `query:"date" validate:"omitempty,datetime=2006-01-02"`
	Timezone string `query:"timezone" validate:"omitempty,timezone"`
}

func (r *GetActivitiesRequest) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}

	if r.Timezone == "" {
		r.Timezone = "Asia/Makassar"
	}
}

type GetActivitiesResponse struct {
	Items []Activity `json:"items"`
	Meta  types.Meta `json:"meta"`
}

type Activity struct {
	Id                  string    `json:"id" db:"id"`
	ServiceAdvisorName  string    `json:"service_advisor_name" db:"service_advisor_name"`
	Status              string    `json:"status" db:"status"`
	TotalPotentialLeads int       `json:"total_potential_leads" db:"total_potential_leads"`
	TotalLeads          int       `json:"total_leads" db:"total_leads"`
	TotalRevenue        float64   `json:"total_revenue" db:"total_revenue"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}
