package ports

import (
	"codebase-app/internal/module/wac/entity"
	"context"
)

var privateFolder = "storage/private"

type WACRepository interface {
	CreateWAC(ctx context.Context, req *entity.CreateWACRequest) error
}

type WACService interface {
	CreateWAC(ctx context.Context, req *entity.CreateWACRequest) error
}
