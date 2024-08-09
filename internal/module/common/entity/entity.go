package entity

type XxxRequest struct {
}

// type GetAreasResponse struct {
// }

type CommonResponse struct {
	Id   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
