package entity

import "codebase-app/pkg/types"

type GetClientsRequest struct {
	Search   string `query:"search" validate:"omitempty,min=3"`
	Page     int    `query:"page" validate:"required"`
	Paginate int    `query:"paginate" validate:"required"`
}

func (r *GetClientsRequest) SetDefault() {
	if r.Paginate < 1 {
		r.Paginate = 10
	}

	if r.Page < 1 {
		r.Page = 1
	}
}

type GetClientsResponse struct {
	Items []Client   `json:"items"`
	Meta  types.Meta `json:"meta"`
}

type Client struct {
	Id          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	VLicense    string `json:"vehicle_license_number" db:"vehicle_license_number"`
	VehicleType string `json:"vehicle_type" db:"vehicle_type"`
	Phone       string `json:"phone" db:"phone"`
}
