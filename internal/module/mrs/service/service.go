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

func (s *mrsService) RenewWAC(ctx context.Context, req *entity.RenewWACRequest) error {
	return s.repo.RenewWAC(ctx, req)
}

func (s *mrsService) DeleteFollowUp(ctx context.Context, req *entity.DeleteFollowUpRequest) error {
	return s.repo.DeleteFollowUp(ctx, req)
}
