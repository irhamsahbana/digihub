package entity

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
	VehicleLicenseNumber string       `json:"vehicle_license_number"`
	WhatsappNumber       string       `json:"whatsapp_number"`
	IsOffered            bool         `json:"is_offered"`
	InvoiceNumber        *string      `json:"invoice_number"`
	Revenue              float64      `json:"revenue"`
	Status               string       `json:"status"`
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
