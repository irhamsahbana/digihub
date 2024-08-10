package entity

import "time"

type CreateWACRequest struct {
	UserId string

	Name                      string             `json:"name" validate:"required"`
	VehicleRegistrationNumber string             `json:"vehicle_registration_number" validate:"min=5,max=20"`
	VehicleTypeId             string             `json:"vehicle_type_id" validate:"ulid,exist=vehicle_types.id"`
	WhatsAppNumber            string             `json:"whatsapp_number" validate:"required"`
	VehicleConditions         []VehicleCondition `json:"vehicle_conditions" validate:"required,dive"`
}

type VehicleCondition struct {
	PotencyId        string  `json:"potency_id" validate:"ulid,exist=potencies.id"`
	AreaId           string  `json:"area_id" validate:"ulid,exist=areas.id"`
	ServiceAdvisorId *string `json:"service_advisor_id" validate:"omitempty,ulid,exist=users.id"`
	Image            string  `json:"image" validate:"base64"`
	Notes            *string `json:"notes"`

	Path string
}

type CreateWACResponse struct {
	Id string `json:"id"`
}

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
