package entity

type CreateWACRequest struct {
	UserId string

	Name                      string             `json:"name" validate:"required"`
	VehicleRegistrationNumber string             `json:"vehicle_registration_number" validate:"min=5,max=10"`
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
