package ports

import (
	"codebase-app/internal/module/common/entity"
	"context"
)

type CommonRepository interface {
	GetAreas(ctx context.Context) ([]entity.AreaResponse, error)
	GetPotencies(ctx context.Context) ([]entity.CommonResponse, error)
	GetVehicleTypes(ctx context.Context) ([]entity.CommonResponse, error)

	GetEmployees(ctx context.Context, req *entity.GetEmployeesRequest) (entity.GetEmployeesResult, error)
	GetBranches(ctx context.Context, req *entity.GetBranchesRequest) (entity.GetBranchesResponse, error)
}

type CommonService interface {
	GetAreas(ctx context.Context) ([]entity.AreaResponse, error)
	GetPotencies(ctx context.Context) ([]entity.CommonResponse, error)
	GetVehicleTypes(ctx context.Context) ([]entity.CommonResponse, error)

	GetEmployees(ctx context.Context, req *entity.GetEmployeesRequest) (entity.GetEmployeesResponse, error)
	GetBranches(ctx context.Context, req *entity.GetBranchesRequest) (entity.GetBranchesResponse, error)
}
