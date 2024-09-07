package ports

import (
	"codebase-app/internal/module/wac/entity"
	"context"
)

type WACRepository interface {
	CreateWAC(ctx context.Context, req *entity.CreateWACRequest) (entity.CreateWACResponse, error)
	GetWACs(ctx context.Context, req *entity.GetWACsRequest) (entity.GetWACsResponse, error)
	GetWAC(ctx context.Context, req *entity.GetWACRequest) (entity.GetWACResponse, error)
	AddRevenue(ctx context.Context, req *entity.AddWACRevenueRequest) error
	AddRevenues(tx context.Context, req *entity.AddWACRevenuesRequest) error

	OfferWAC(ctx context.Context, req *entity.OfferWACRequest) (entity.OfferWACResponse, error)
	IsWACCreator(ctx context.Context, userId, WACId string) (bool, error)
	IsWACStatus(ctx context.Context, WACId, status string) (bool, error)
}

type WACService interface {
	CreateWAC(ctx context.Context, req *entity.CreateWACRequest) (entity.CreateWACResponse, error)
	GetWACs(ctx context.Context, req *entity.GetWACsRequest) (entity.GetWACsResponse, error)
	GetWAC(ctx context.Context, req *entity.GetWACRequest) (entity.GetWACResponse, error)

	OfferWAC(ctx context.Context, req *entity.OfferWACRequest) (entity.OfferWACResponse, error)
	AddRevenue(ctx context.Context, req *entity.AddWACRevenueRequest) (entity.AddWACRevenueResponse, error)
	AddRevenues(tx context.Context, req *entity.AddWACRevenuesRequest) (entity.AddWACRevenueResponse, error)
}
