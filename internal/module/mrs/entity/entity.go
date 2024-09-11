package entity

import (
	"codebase-app/pkg/types"
	"time"
)

type GetMRSsRequest struct {
	UserId string

	Page     int `query:"page" validate:"required"`
	Paginate int `query:"paginate" validate:"required"`
}

func (r *GetMRSsRequest) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type GetMRSsResponse struct {
	Items []MRSItem  `json:"items"`
	Meta  types.Meta `json:"meta"`
}

type MRSItem struct {
	Id             string     `json:"id" db:"id"`
	Client         string     `json:"client" db:"client"`
	ServiceAdvisor string     `json:"service_advisor" db:"service_advisor"`
	FollowUpAt     *time.Time `json:"follow_up_at" db:"follow_up_at"`
}

type RenewWACRequest struct {
	UserId              string
	WacId               string   `params:"id" validate:"ulid"`
	VehicleConditionIds []string `json:"vehicle_condition_ids" validate:"required,dive,ulid"`
}
