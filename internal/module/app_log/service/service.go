package service

import (
	"codebase-app/internal/module/app_log/entity"
	"codebase-app/internal/module/app_log/ports"
	"context"
)

var _ ports.AppLogService = &appLogService{}

type appLogService struct {
	repo ports.AppLogRepository
}

func NewAppLogService(repo ports.AppLogRepository) *appLogService {
	return &appLogService{
		repo: repo,
	}
}

func (s *appLogService) CreateLog(ctx context.Context, req *entity.CreateLogRequest) error {
	return s.repo.CreateLog(ctx, req)
}
