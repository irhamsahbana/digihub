package entity

type CreatePromotionRequest struct {
	Title string  `json:"title" db:"title" validate:"required,max=255"`
	Image string  `json:"image" validate:"base64"`
	Link  *string `json:"link" db:"link" validate:"omitempty,url"`

	Path string `db:"path"`
}

func (r *CreatePromotionRequest) RemoveImage() {
	r.Image = ""
}

type DeletePromotionRequest struct {
	Id string `json:"id" validate:"ulid"`

	Path string `db:"path"`
}

type Promotion struct {
	Id    string  `json:"id" db:"id"`
	Title string  `json:"title" db:"title"`
	Image string  `json:"image"`
	Link  *string `json:"link" db:"link"`
	Path  string  `json:"-" db:"path"` // json:"-" to hide the field in response
}
