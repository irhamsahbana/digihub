package entity

type MarkWIPRequest struct {
	UserId string

	Id string `params:"id" validate:"ulid"`
}

type MarkWIPResponse struct {
	Id string `json:"id"`
}
