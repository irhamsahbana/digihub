package entity

import "time"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	ExpiredAt   time.Time `json:"expired_at"`
}

type CommonResponse struct {
	Id   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type GetProfileRequest struct {
	UserId string
}

type GetProfileResponse struct {
	Id      string         `json:"id" db:"id"`
	Name    string         `json:"name" db:"name"`
	Email   string         `json:"email" db:"email"`
	WANum   string         `json:"whatsapp_number" db:"whatsapp_number"`
	Image   *string        `json:"image"`
	Role    CommonResponse `json:"role"`
	Branch  CommonResponse `json:"branch"`
	Section CommonResponse `json:"section"`
}
