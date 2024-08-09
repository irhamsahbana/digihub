package entity

type GetEmployeesRequest struct {
	Page     int    `query:"page" validate:"required"`
	Paginate int    `query:"paginate" validate:"required"`
	Role     string `query:"role" validate:"omitempty,oneof=service_advisor technician"`
	BranchId string `query:"branch_id" validate:"omitempty,exist=branches.id"`
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
