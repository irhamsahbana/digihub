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

func (s *commonService) GetBranches(ctx context.Context, req *entity.GetBranchesRequest) (entity.GetBranchesResponse, error) {
	return s.repo.GetBranches(ctx, req)
}

func (s *commonService) GetAreas(ctx context.Context) ([]entity.AreaResponse, error) {
	return s.repo.GetAreas(ctx)
}

func (s *commonService) GetPotencies(ctx context.Context, req *entity.GetPotenciesRequest) ([]entity.GetPotencyResponse, error) {
	return s.repo.GetPotencies(ctx, req)
}

func (s *commonService) GetVehicleTypes(ctx context.Context) ([]entity.CommonResponse, error) {
	return s.repo.GetVehicleTypes(ctx)
}

func (s *commonService) GetEmployees(ctx context.Context, req *entity.GetEmployeesRequest) (entity.GetEmployeesResponse, error) {
	var (
		resp entity.GetEmployeesResponse
	)
	resp.Items = make([]entity.EmployeeItem, 0, req.Paginate)

	result, err := s.repo.GetEmployees(ctx, req)
	if err != nil {
		return resp, err
	}

	resp.Meta = result.Meta

	for _, item := range result.Items {
		resp.Items = append(resp.Items, entity.EmployeeItem{
			UserId:      item.UserId,
			Name:        item.Name,
			Email:       item.Email,
			WhatsappNum: item.WhatsappNum,
			EIBranch: entity.CommonResponse{
				Id:   item.BranchId,
				Name: item.BranchName,
			},
			EISection: entity.CommonResponse{
				Id:   item.SectionId,
				Name: item.SectionName,
			},
			EIRole: entity.CommonResponse{
				Id:   item.RoleId,
				Name: item.RoleName,
			},
		})
	}

	return resp, nil
}

func (s *commonService) GetRoles(ctx context.Context) ([]entity.CommonResponse, error) {
	return s.repo.GetRoles(ctx)
}

func (s *commonService) GetHTIBrands(ctx context.Context) ([]entity.CommonResponse, error) {
	return s.repo.GetHTIBrands(ctx)
}

func (s *commonService) GetHTIModels(ctx context.Context, req *entity.GetHTIModelsRequest) ([]entity.CommonResponse, error) {
	return s.repo.GetHTIModels(ctx, req)
}

func (s *commonService) GetHTITypes(ctx context.Context, req *entity.GetHTITypesRequest) ([]entity.CommonResponse, error) {
	return s.repo.GetHTITypes(ctx, req)
}

func (s *commonService) GetHTIYears(ctx context.Context, req *entity.GetHTIYearsRequest) ([]entity.CommonResponse, error) {
	return s.repo.GetHTIYears(ctx, req)
}

func (s *commonService) GetHTIPurchase(ctx context.Context, req *entity.GetHTIPurchaseRequest) (entity.GetHTIPurchaseResponse, error) {
	return s.repo.GetHTIPurchase(ctx, req)
}

func (s *commonService) GetHTIValuations(ctx context.Context, req *entity.GetHTIValuationsRequest) (entity.GetHTIValuationsResponse, error) {
	return s.repo.GetHTIValuations(ctx, req)
}
