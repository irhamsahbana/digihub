package ports

import (
	"codebase-app/internal/module/mrs/entity"
	"context"
)

type MRSRepository interface {
	GetMRSs(ctx context.Context, req *entity.GetMRSsRequest) (entity.GetMRSsResponse, error)
	RenewWAC(ctx context.Context, req *entity.RenewWACRequest) error
}

type MRSService interface {
	GetMRSs(ctx context.Context, req *entity.GetMRSsRequest) (entity.GetMRSsResponse, error)
	RenewWAC(ctx context.Context, req *entity.RenewWACRequest) error
}
