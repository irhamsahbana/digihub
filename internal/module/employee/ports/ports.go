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
}

type EmployeeService interface {
	GetEmployee(ctx context.Context, req *entity.GetEmployeeRequest) (entity.GetEmployeeResponse, error)
	UpdateEmployee(ctx context.Context, req *entity.UpdateEmployeeRequest) error
	CreateEmployee(ctx context.Context, req *entity.CreateEmployeeRequest) error
	DeleteEmployee(ctx context.Context, req *entity.DeleteEmployeeRequest) error
}
