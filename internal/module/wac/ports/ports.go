package ports

import (
	"codebase-app/internal/module/wac/entity"
	"context"
)

type WACRepository interface {
	CreateWAC(ctx context.Context, req *entity.CreateWACRequest) (entity.CreateWACResponse, error)
	GetWACs(ctx context.Context, req *entity.GetWACsRequest) (entity.GetWACsResponse, error)
	GetWAC(ctx context.Context, req *entity.GetWACRequest) (entity.GetWACResponse, error)

	OfferWAC(ctx context.Context, req *entity.OfferWACRequest) (entity.OfferWACResponse, error)
	IsWACCreator(ctx context.Context, userId, WACId string) (bool, error)
	IsWACOffered(ctx context.Context, WACId string) (bool, error)
}

type WACService interface {
	CreateWAC(ctx context.Context, req *entity.CreateWACRequest) (entity.CreateWACResponse, error)
	GetWACs(ctx context.Context, req *entity.GetWACsRequest) (entity.GetWACsResponse, error)
	GetWAC(ctx context.Context, req *entity.GetWACRequest) (entity.GetWACResponse, error)

	OfferWAC(ctx context.Context, req *entity.OfferWACRequest) (entity.OfferWACResponse, error)
}
