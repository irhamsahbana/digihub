package entity

type CreateWACRequest struct {
	Name                      string             `json:"name" validate:"required"`
	VehicleRegistrationNumber string             `json:"vehicle_registration_number" validate:"required"`
	VehicleTypeId             string             `json:"vehicle_type_id" validate:"required,ulid"`
	WhatsAppNumber            string             `json:"whatsapp_number" validate:"required"`
	VehicleConditions         []VehicleCondition `json:"vehicle_conditions" validate:"required,dive"`
}

type VehicleCondition struct {
	PotencyId string `json:"potency_id" validate:"required,ulid"`
	AreaId    string `json:"area_id" validate:"required,ulid"`
	Image     string `json:"image" validate:"required,base64"`
}

type XxxResponse struct {
}

type XxxResult struct {
}
