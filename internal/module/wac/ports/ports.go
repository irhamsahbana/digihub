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
	MarkWIP(ctx context.Context, req *entity.MarkWIPRequest) error

	OfferWAC(ctx context.Context, req *entity.OfferWACRequest) (entity.OfferWACResponse, error)
	IsWACCreator(ctx context.Context, userId, WACId string) (bool, error)
	IsWACStatus(ctx context.Context, WACId, status string) (bool, error)
}

type WACService interface {
	CreateWAC(ctx context.Context, req *entity.CreateWACRequest) (entity.CreateWACResponse, error)
	GetWACs(ctx context.Context, req *entity.GetWACsRequest) (entity.GetWACsResponse, error)
	GetWAC(ctx context.Context, req *entity.GetWACRequest) (entity.GetWACResponse, error)
	MarkWIP(ctx context.Context, req *entity.MarkWIPRequest) (entity.MarkWIPResponse, error)

	OfferWAC(ctx context.Context, req *entity.OfferWACRequest) (entity.OfferWACResponse, error)
	AddRevenue(ctx context.Context, req *entity.AddWACRevenueRequest) (entity.AddWACRevenueResponse, error)
}
