package ports

import (
	"codebase-app/internal/module/dashboard/entity"
	"context"
)

type DashboardRepository interface {
	GetLeadsTrends(ctx context.Context, request *entity.LeadTrendsRequest) ([]entity.LeadTrendsResponse, error)
	GetWACSummary(ctx context.Context, request *entity.WACSummaryRequest) (entity.WACSummaryResponse, error)
}

type DashboardService interface {
	GetLeadsTrends(ctx context.Context, request *entity.LeadTrendsRequest) ([]entity.LeadTrendsResponse, error)
	GetWACSummary(ctx context.Context, request *entity.WACSummaryRequest) (entity.WACSummaryResponse, error)
}
