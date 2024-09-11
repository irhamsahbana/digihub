package service

import (
	"codebase-app/internal/module/mrs/entity"
	"codebase-app/internal/module/mrs/ports"
	"context"
)

var _ ports.MRSService = &mrsService{}

type mrsService struct {
	repo ports.MRSRepository
}

func NewMRSService(repo ports.MRSRepository) *mrsService {
	return &mrsService{
		repo: repo,
	}
}

func (s *mrsService) GetMRSs(ctx context.Context, req *entity.GetMRSsRequest) (entity.GetMRSsResponse, error) {
	return s.repo.GetMRSs(ctx, req)
}
