package entity

type GetBranchesRequest struct {
	Page     int `query:"page" validate:"required"`
	Paginate int `query:"paginate" validate:"required"`
}

func (r *GetBranchesRequest) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type GetBranchesResponse struct {
	Items []CommonResponse `json:"items"`
	Meta  Meta             `json:"meta"`
}

type CreateBranchRequest struct {
	Name    string  `json:"name" validate:"required"`
	Address *string `json:"address"`
}

type UpdateBranchRequest struct {
	Id      string  `json:"id" validate:"required,ulid"`
	Name    string  `json:"name" validate:"required"`
	Address *string `json:"address"`
}
