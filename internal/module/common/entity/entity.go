package entity

type GetEmployeesRequest struct {
	UserId string

	Page      int    `query:"page" validate:"required"`
	Paginate  int    `query:"paginate" validate:"required"`
	Role      string `query:"role" validate:"omitempty,oneof=service_advisor technician"`
	BranchId  string `query:"branch_id" validate:"omitempty,exist=branches.id"`
	SectionId string `query:"section_id" validate:"omitempty,exist=potencies.id"`
	IncludeMe bool   `query:"include_me"`
	Search    string `query:"search" validate:"omitempty,min=3"`
}

func (r *GetEmployeesRequest) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type GetEmployeesResponse struct {
	Items []EmployeeItem `json:"items"`
	Meta  Meta           `json:"meta"`
}

type EmployeeItem struct {
	UserId      string         `json:"user_id"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	WhatsappNum string         `json:"whatsapp_number"`
	EIBranch    CommonResponse `json:"branch"`
	EISection   CommonResponse `json:"section"`
	EIRole      CommonResponse `json:"role"`
}

type CommonResponse struct {
	Id   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type AreaResponse struct {
	CommonResponse
	Type string `json:"type" db:"type"`
}

/*
	Results
*/

type GetEmployeesResult struct {
	Items []EmployeeResult
	Meta  Meta
}

type EmployeeResult struct {
	UserId      string `db:"user_id"`
	BranchId    string `db:"branch_id"`
	SectionId   string `db:"section_id"`
	RoleId      string `db:"role_id"`
	Name        string `db:"name"`
	Email       string `db:"email"`
	WhatsappNum string `db:"whatsapp_number"`
	BranchName  string `db:"branch_name"`
	SectionName string `db:"section_name"`
	RoleName    string `db:"role_name"`
}

type GetProfileResponse struct {
	Id     string         `json:"id" db:"id"`
	Branch CommonResponse `json:"branch"`
}

type GetPotenciesRequest struct {
	UserId string
}

type GetPotencyResponse struct {
	Id     string          `json:"id" db:"id"`
	Name   string          `json:"name" db:"name"`
	User   *CommonResponse `json:"user"`
	Branch *CommonResponse `json:"branch"`
}

type GetHTIValuationsRequest struct {
	Brand    string `query:"brand" validate:"required"`
	Model    string `query:"model" validate:"required"`
	Page     int    `query:"page" validate:"required"`
	Paginate int    `query:"paginate" validate:"required"`
}

func (r *GetHTIValuationsRequest) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type GetHTIValuationsResponse struct {
	Items []HTIValuationItem `json:"items"`
	Meta  Meta               `json:"meta"`
}

type HTIValuationItem struct {
	Type        string  `json:"type" db:"type"`
	Year        int     `json:"year" db:"year"`
	MinPurchase float64 `json:"min_purchase" db:"min_purchase"`
	MaxPurchase float64 `json:"max_purchase" db:"max_purchase"`
}
