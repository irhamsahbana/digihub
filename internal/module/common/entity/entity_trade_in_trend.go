package entity

type GetHTIBrandsRequest struct {
}

type GetHTIBrandsResponse struct {
	Items []CommonResponse `json:"items"`
}

type GetHTIModelsRequest struct {
	Brand string `query:"brand" validate:"required"`
}

type GetHTIModelsResponse struct {
	Items []CommonResponse `json:"items"`
}

type GetHTITypesRequest struct {
	Brand string `query:"brand" validate:"required"`
	Model string `query:"model" validate:"required"`
}

type GetHTITypesResponse struct {
	Items []CommonResponse `json:"items"`
}

type GetHTIYearsRequest struct {
	Brand string `query:"brand" validate:"required"`
	Model string `query:"model" validate:"required"`
	Type  string `query:"type" validate:"required"`
}

type GetHTIYearsResponse struct {
	Items []CommonResponse `json:"items"`
}

type GetHTIPurchaseRequest struct {
	Brand string `query:"brand" validate:"required"`
	Model string `query:"model" validate:"required"`
	Type  string `query:"type" validate:"required"`
	Year  string `query:"year" validate:"required,numeric"`
}

type GetHTIPurchaseResponse struct {
	MinPurchase int `db:"min_purchase" json:"min_purchase"`
	MaxPurchase int `db:"max_purchase" json:"max_purchase"`
}
