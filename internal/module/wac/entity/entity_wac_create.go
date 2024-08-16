package entity

type CreateWACRequest struct {
	UserId string `validate:"required,ulid,exist=users.id"`

	Name                      string             `json:"name" validate:"required"`
	VehicleRegistrationNumber string             `json:"vehicle_registration_number" validate:"min=3,max=15"`
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

func (r *CreateWACRequest) RemoveBase64() {
	for i := range r.VehicleConditions {
		r.VehicleConditions[i].Image = ""
	}
}

type CreateWACResponse struct {
	Id string `json:"id"`
}
