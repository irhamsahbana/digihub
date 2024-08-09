package service

import (
	"codebase-app/internal/module/common/entity"
	"codebase-app/internal/module/common/ports"
	"context"
)

var _ ports.CommonService = &commonService{}

type commonService struct {
	repo ports.CommonRepository
}

func NewCommonService(repo ports.CommonRepository) *commonService {
	return &commonService{
		repo: repo,
	}
}

func (s *commonService) GetAreas(ctx context.Context) ([]entity.CommonResponse, error) {
	return s.repo.GetAreas(ctx)
}

func (s *commonService) GetPotencies(ctx context.Context) ([]entity.CommonResponse, error) {
	return s.repo.GetPotencies(ctx)
}

func (s *commonService) GetVehicleTypes(ctx context.Context) ([]entity.CommonResponse, error) {
	return s.repo.GetVehicleTypes(ctx)
}
