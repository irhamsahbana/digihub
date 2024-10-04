package ports

import (
	"codebase-app/internal/module/common/entity"
	"context"
)

type CommonRepository interface {
	GetAreas(ctx context.Context) ([]entity.AreaResponse, error)
	GetPotencies(ctx context.Context, req *entity.GetPotenciesRequest) ([]entity.GetPotencyResponse, error)
	GetVehicleTypes(ctx context.Context) ([]entity.CommonResponse, error)

	GetEmployees(ctx context.Context, req *entity.GetEmployeesRequest) (entity.GetEmployeesResult, error)
	GetBranches(ctx context.Context, req *entity.GetBranchesRequest) (entity.GetBranchesResponse, error)
	GetRoles(ctx context.Context) ([]entity.CommonResponse, error)

	GetHTIBrands(ctx context.Context) ([]entity.CommonResponse, error)
	GetHTIModels(ctx context.Context, req *entity.GetHTIModelsRequest) ([]entity.CommonResponse, error)
	GetHTITypes(ctx context.Context, req *entity.GetHTITypesRequest) ([]entity.CommonResponse, error)
	GetHTIYears(ctx context.Context, req *entity.GetHTIYearsRequest) ([]entity.CommonResponse, error)
	GetHTIPurchase(ctx context.Context, req *entity.GetHTIPurchaseRequest) (entity.GetHTIPurchaseResponse, error)
	GetHTIValuations(ctx context.Context, req *entity.GetHTIValuationsRequest) (entity.GetHTIValuationsResponse, error)
}

type CommonService interface {
	GetAreas(ctx context.Context) ([]entity.AreaResponse, error)
	GetPotencies(ctx context.Context, req *entity.GetPotenciesRequest) ([]entity.GetPotencyResponse, error)
	GetVehicleTypes(ctx context.Context) ([]entity.CommonResponse, error)

	GetEmployees(ctx context.Context, req *entity.GetEmployeesRequest) (entity.GetEmployeesResponse, error)
	GetBranches(ctx context.Context, req *entity.GetBranchesRequest) (entity.GetBranchesResponse, error)
	GetRoles(ctx context.Context) ([]entity.CommonResponse, error)

	GetHTIBrands(ctx context.Context) ([]entity.CommonResponse, error)
	GetHTIModels(ctx context.Context, req *entity.GetHTIModelsRequest) ([]entity.CommonResponse, error)
	GetHTITypes(ctx context.Context, req *entity.GetHTITypesRequest) ([]entity.CommonResponse, error)
	GetHTIYears(ctx context.Context, req *entity.GetHTIYearsRequest) ([]entity.CommonResponse, error)
	GetHTIPurchase(ctx context.Context, req *entity.GetHTIPurchaseRequest) (entity.GetHTIPurchaseResponse, error)
	GetHTIValuations(ctx context.Context, req *entity.GetHTIValuationsRequest) (entity.GetHTIValuationsResponse, error)
}
