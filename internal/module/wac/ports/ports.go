package ports

import (
	"codebase-app/internal/module/wac/entity"
	"context"
)

var privateFolder = "storage/private"

type WACRepository interface {
	CreateWAC(ctx context.Context, req *entity.CreateWACRequest) (entity.CreateWACResponse, error)
	GetWACs(ctx context.Context, req *entity.GetWACsRequest) (entity.GetWACsResponse, error)
}

type WACService interface {
	CreateWAC(ctx context.Context, req *entity.CreateWACRequest) (entity.CreateWACResponse, error)
	GetWACs(ctx context.Context, req *entity.GetWACsRequest) (entity.GetWACsResponse, error)
}
