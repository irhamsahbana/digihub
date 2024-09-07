package entity

type AddWACRevenueRequest struct {
	UserId string `validate:"ulid"`

	Id            string  `params:"id" validate:"ulid,exist=walk_around_checks.id"`
	InvoiceNumber string  `json:"invoice_number" validate:"required,min=4,max=255"`
	TotalRevenue  float64 `json:"total_revenue" validate:"numeric"`
}

type AddWACRevenueResponse struct {
	Id string `json:"id" db:"id"`
}

type AddWACRevenuesRequest struct {
	UserId string `validate:"required,ulid"`

	Id       string       `params:"id" validate:"required,ulid"`
	Revenues []WACRevenue `prop:"revenues" validate:"required,dive"`
}

type WACRevenue struct {
	VehicleConditionId string  `json:"vehicle_condition_id" validate:"required,ulid"`
	InvoiceNumber      string  `json:"invoice_number" validate:"required,min=4,max=255"`
	Revenue            float64 `json:"revenue" validate:"numeric"`
}
