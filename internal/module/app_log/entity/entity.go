package entity

type CreateLogRequest struct {
	AppId    string  `json:"app_id" validate:"required"`
	LogLevel string  `json:"log_level" validate:"required"`
	Info     *string `json:"info"`
	Message  string  `json:"message" validate:"required"`
}
