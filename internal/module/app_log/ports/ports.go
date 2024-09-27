package ports

import (
	"codebase-app/internal/module/app_log/entity"
	"context"
)

type AppLogRepository interface {
	CreateLog(ctx context.Context, req *entity.CreateLogRequest) error
}

type AppLogService interface {
	CreateLog(ctx context.Context, req *entity.CreateLogRequest) error
}
