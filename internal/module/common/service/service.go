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
func (s *commonService) GetEmployees(ctx context.Context, req *entity.GetEmployeesRequest) (entity.GetEmployeesResponse, error) {
	var (
		resp entity.GetEmployeesResponse
	)

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
