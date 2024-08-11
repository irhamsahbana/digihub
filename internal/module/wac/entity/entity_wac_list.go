package entity

import "time"

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
	Meta  Meta                 `json:"meta"`
}

type WacItem struct {
	Id         string    `json:"id" db:"id"`
	ClientName string    `json:"client_name" db:"client_name"`
	Status     string    `json:"status" db:"status"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
