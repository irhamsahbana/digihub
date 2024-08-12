package entity

type AddWACRevenueRequest struct {
	UserId string `validate:"ulid"`

	Id            string  `json:"id" validate:"ulid,exist=walk_around_checks.id"`
	InvoiceNumber string  `json:"invoice_number" validate:"required,min=4,max=255"`
	TotalRevenue  float64 `json:"total_revenue" validate:"numeric"`
}

type AddWACRevenueResponse struct {
	Id string `json:"id" db:"id"`
}
