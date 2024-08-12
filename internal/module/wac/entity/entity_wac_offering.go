package entity

type OfferWACRequest struct {
	UserId string `validate:"ulid"`

	Id          string           `params:"id" validate:"ulid,exist=walk_around_checks.id"`
	IsUsedCar   bool             `json:"is_used_car"`
	VConditions []OfferCondition `json:"vehicle_conditions" validate:"omitempty,min=1,dive"`
}

type OfferCondition struct {
	Id           string  `json:"id" validate:"ulid"`
	IsInterested bool    `json:"is_interested"`
	Notes        *string `json:"notes" validate:"omitempty,max=255"`
}

type OfferWACResponse struct {
	Id string `json:"id" db:"id"`
}
