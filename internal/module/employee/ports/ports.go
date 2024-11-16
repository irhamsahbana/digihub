package ports

import (
	"codebase-app/internal/module/employee/entity"
	"context"
)

type EmployeeRepository interface {
	GetEmployee(ctx context.Context, req *entity.GetEmployeeRequest) (entity.GetEmployeeResponse, error)
	UpdateEmployee(ctx context.Context, req *entity.UpdateEmployeeRequest) error
	CreateEmployee(ctx context.Context, req *entity.CreateEmployeeRequest) error
	DeleteEmployee(ctx context.Context, req *entity.DeleteEmployeeRequest) error

	ImportEmployees(ctx context.Context, data []entity.ImportEmployeeRow) error

	GetBranches(ctx context.Context) ([]entity.Common, error)
	GetPotencies(ctx context.Context) ([]entity.Common, error)
	GetRoles(ctx context.Context) ([]entity.Common, error)

	IsEmailExist(ctx context.Context, email string) error
}

type EmployeeService interface {
	GetEmployee(ctx context.Context, req *entity.GetEmployeeRequest) (entity.GetEmployeeResponse, error)
	UpdateEmployee(ctx context.Context, req *entity.UpdateEmployeeRequest) error
	CreateEmployee(ctx context.Context, req *entity.CreateEmployeeRequest) error
	DeleteEmployee(ctx context.Context, req *entity.DeleteEmployeeRequest) error

	ImportEmployees(ctx context.Context, req *entity.ImportEmployeesRequest) error
}
