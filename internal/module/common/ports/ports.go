package ports

import (
	"codebase-app/internal/module/common/entity"
	"context"
)

type CommonRepository interface {
	GetAreas(ctx context.Context) ([]entity.CommonResponse, error)
	GetPotencies(ctx context.Context) ([]entity.CommonResponse, error)
	GetVehicleTypes(ctx context.Context) ([]entity.CommonResponse, error)
}

type CommonService interface {
	GetAreas(ctx context.Context) ([]entity.CommonResponse, error)
	GetPotencies(ctx context.Context) ([]entity.CommonResponse, error)
	GetVehicleTypes(ctx context.Context) ([]entity.CommonResponse, error)
}
