package ports

import (
	"codebase-app/internal/module/dashboard/entity"
	"context"
)

type DashboardRepository interface {
	GetLeadsTrends(ctx context.Context, request *entity.LeadTrendsRequest) ([]entity.LeadTrendsResponse, error)
	GetWACSummary(ctx context.Context, request *entity.WACSummaryRequest) (entity.WACSummaryResponse, error)
	GetWACSummaryTechnician(ctx context.Context, request *entity.WACSummaryRequest) (entity.TechWACSummaryResponse, error)

	GetWACLineChart(ctx context.Context, request *entity.GetWACLineChartRequest) ([]entity.GetWACLineChartResponse, error)
	GetActivities(ctx context.Context, request *entity.GetActivitiesRequest) (entity.GetActivitiesResponse, error)
}

type DashboardService interface {
	GetLeadsTrends(ctx context.Context, request *entity.LeadTrendsRequest) ([]entity.LeadTrendsResponse, error)
	GetWACSummary(ctx context.Context, request *entity.WACSummaryRequest) (entity.WACSummaryResponse, error)
	GetWACSummaryTechnician(ctx context.Context, request *entity.WACSummaryRequest) (entity.TechWACSummaryResponse, error)

	GetWACLineChart(ctx context.Context, request *entity.GetWACLineChartRequest) ([]entity.GetWACLineChartResponse, error)
	GetActivities(ctx context.Context, request *entity.GetActivitiesRequest) (entity.GetActivitiesResponse, error)
}
