package entity

type GetEmployeeRequest struct {
	Id string `json:"id" validate:"ulid"`
}

type GetEmployeeResponse struct {
	Id       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	WANumber string `json:"whatssap_number" db:"whatsapp_number"`
	Branch   Common `json:"branch"`
	Section  Common `json:"section"`
	Role     Common `json:"role"`

	BranchId    string `json:"-" db:"branch_id"`
	SectionId   string `json:"-" db:"section_id"`
	RoleId      string `json:"-" db:"role_id"`
	BranchName  string `json:"-" db:"branch_name"`
	SectionName string `json:"-" db:"section_name"`
	RoleName    string `json:"-" db:"role_name"`
}

type Common struct {
	Id   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type UpdateEmployeeRequest struct {
	Id        string `params:"id" validate:"ulid"`
	BranchId  string `json:"branch_id" validate:"ulid,exist=branches.id"`
	SectionId string `json:"section_id" validate:"ulid,exist=potencies.id"`
	RoleId    string `json:"role_id" validate:"ulid,exist=roles.id"`
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"email"`
	WANumber  string `json:"whatsapp_number"`
	Password  string `json:"password" validate:"omitempty,strong_password"`
	PassConf  string `json:"password_confirmation" validate:"eqfield=Password"`
}

type CreateEmployeeRequest struct {
	BranchId  string `json:"branch_id" validate:"ulid,exist=branches.id"`
	SectionId string `json:"section_id" validate:"ulid,exist=potencies.id"`
	RoleId    string `json:"role_id" validate:"ulid,exist=roles.id"`
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"email"`
	WANumber  string `json:"whatsapp_number"`
	Password  string `json:"password" validate:"required,strong_password"`
	PassConf  string `json:"password_confirmation" validate:"eqfield=Password"`
}

type DeleteEmployeeRequest struct {
	Id string `json:"id" validate:"ulid"`
}
