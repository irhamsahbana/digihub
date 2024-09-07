package entity

import (
	"codebase-app/pkg/types"
	"time"
)

type GetWACsRequest struct {
	UserId string

	Page     int    `query:"page" validate:"required"`
	Paginate int    `query:"paginate" validate:"required"`
	Query    string `query:"query" validate:"omitempty,min=3"`
	Status   string `query:"status" validate:"omitempty,oneof=created offered wip completed"`
}

func (r *GetWACsRequest) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type GetWACsResponse struct {
	Items map[string][]WacItem `json:"items"`
	Meta  types.Meta           `json:"meta"`
}

type WacItem struct {
	Id                  string    `json:"id" db:"id"`
	ClientName          string    `json:"client_name" db:"client_name"`
	TotalPotentialLeads int       `json:"total_potential_leads" db:"total_potential_leads"`
	TotalLeads          int       `json:"total_leads" db:"total_leads"`
	TotalLeadsCompleted int       `json:"total_leads_completed" db:"total_leads_completed"`
	TotalFollowUps      int       `json:"total_follow_ups" db:"total_follow_ups"`
	Status              string    `json:"status" db:"status"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}
