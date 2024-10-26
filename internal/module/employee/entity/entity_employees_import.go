package entity

type ImportEmployeesRequest struct {
	File string `json:"file" validate:"required,base64"`
}

type ImportEmployeeRow struct {
	Name           string `json:"name"`
	BranchName     string `json:"branch_name"`
	SectionName    string `json:"section_name"`
	RoleName       string `json:"role_name"`
	Email          string `json:"email"`
	PasswordHashed string `json:"password"`

	BranchId  string `json:"branch_id"`
	SectionId string `json:"section_id"`
	RoleId    string `json:"role_id"`
}
