package entity

import "time"

type Common struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Area struct {
	Common
	Type string `json:"type"`
}

type GetWACRequest struct {
	Id string `param:"id" validate:"required,ulid"`
}

type GetWACResponse struct {
	Id                   string       `json:"id"`
	ClientName           string       `json:"client_name"`
	SAName               string       `json:"service_advisor_name"`
	BranchName           string       `json:"branch_name"`
	VehicleLicenseNumber string       `json:"vehicle_license_number"`
	WhatsappNumber       string       `json:"whatsapp_number"`
	IsUsedCar            bool         `json:"is_used_car"`
	IsOffered            bool         `json:"is_offered"`
	InvoiceNumber        *string      `json:"invoice_number"`
	Revenue              float64      `json:"revenue"`
	Status               string       `json:"status"`
	CreatedAt            time.Time    `json:"created_at"`
	VehicleType          Common       `json:"vehicle_type"`
	VehicleConditions    []VCondition `json:"vehicle_conditions"`
}

type VCondition struct {
	Id           string        `json:"id"`
	InvoiceNum   *string       `json:"invoice_number"`
	Revenue      float64       `json:"revenue"`
	Potency      Common        `json:"potency"`
	Area         Area          `json:"area"`
	Assginee     AssigneedUser `json:"assigneed_user"`
	Image        string        `json:"image"`
	IsInterested bool          `json:"is_interested"`
	Notes        *string       `json:"notes"`
}

type AssigneedUser struct {
	Branch Common `json:"branch"`
	User   Common `json:"user"`
}

type GetWACPDFLinkRequest struct {
	Id string `params:"id" validate:"required,ulid"`
}

type GetWACPDFLinkResponse struct {
	Link      string `json:"link"`
	Expires   int64  `json:"expires"`
	Signature string `json:"signature"`
}
