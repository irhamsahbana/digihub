package ports

import (
	"codebase-app/internal/module/dashboard/entity"
	"context"
)

type DashboardRepository interface {
	GetLeadsTrends(ctx context.Context, request *entity.LeadTrendsRequest) ([]entity.LeadTrendsResponse, error)
	GetWACSummary(ctx context.Context, request *entity.WACSummaryRequest) (entity.WACSummaryResponse, error)
	GetWACSummaryTechnician(ctx context.Context, request *entity.WACSummaryRequest) (entity.TechWACSummaryResponse, error)
}

type DashboardService interface {
	GetLeadsTrends(ctx context.Context, request *entity.LeadTrendsRequest) ([]entity.LeadTrendsResponse, error)
	GetWACSummary(ctx context.Context, request *entity.WACSummaryRequest) (entity.WACSummaryResponse, error)
	GetWACSummaryTechnician(ctx context.Context, request *entity.WACSummaryRequest) (entity.TechWACSummaryResponse, error)
}
