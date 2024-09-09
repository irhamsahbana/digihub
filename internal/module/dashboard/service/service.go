package service

import (
	"codebase-app/internal/module/dashboard/entity"
	"codebase-app/internal/module/dashboard/ports"
	"context"
)

var _ ports.DashboardService = &DashbaordService{}

type DashbaordService struct {
	repo ports.DashboardRepository
}

func NewDashboardService(repo ports.DashboardRepository) *DashbaordService {
	return &DashbaordService{
		repo: repo,
	}
}

func (s *DashbaordService) GetLeadsTrends(ctx context.Context, request *entity.LeadTrendsRequest) ([]entity.LeadTrendsResponse, error) {
	return s.repo.GetLeadsTrends(ctx, request)
}

func (s *DashbaordService) GetWACSummary(ctx context.Context, request *entity.WACSummaryRequest) (entity.WACSummaryResponse, error) {
	return s.repo.GetWACSummary(ctx, request)
}

func (s *DashbaordService) GetWACSummaryTechnician(ctx context.Context, request *entity.WACSummaryRequest) (entity.TechWACSummaryResponse, error) {
	return s.repo.GetWACSummaryTechnician(ctx, request)
}
