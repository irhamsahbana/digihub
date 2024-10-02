package entity

import "github.com/LukaGiorgadze/gonull"

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
	Id string `params:"id" validate:"ulid"`

	Path string `db:"path"`
}

type UpdatePromotionRequest struct {
	Id string `params:"id" validate:"ulid"`

	Title gonull.Nullable[string] `json:"title"`
	Image gonull.Nullable[string] `json:"image"`
	Link  gonull.Nullable[string] `json:"link"`

	TitleVal string `validate:"omitempty,max=255" prop:"title"`
	ImageVal string `validate:"omitempty,base64" prop:"image"`
	LinkVal  string `validate:"omitempty,url" prop:"link"`

	Path string `db:"path"`
}

func (r *UpdatePromotionRequest) RemoveImage() {
	r.Image = gonull.NewNullable("")
}

func (r *UpdatePromotionRequest) SetValues() {
	r.TitleVal = r.Title.Val
	r.ImageVal = r.Image.Val
	r.LinkVal = r.Link.Val
}

type Promotion struct {
	Id    string  `json:"id" db:"id"`
	Title string  `json:"title" db:"title"`
	Image string  `json:"image"`
	Link  *string `json:"link" db:"link"`

	Path string `json:"-" db:"path"` // json:"-" to hide the field in response
}
