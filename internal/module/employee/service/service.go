package service

import (
	"codebase-app/internal/module/employee/entity"
	"codebase-app/internal/module/employee/ports"
	"context"
)

var _ ports.EmployeeService = &employeeService{}

type employeeService struct {
	repo ports.EmployeeRepository
}

func NewEmployeeService(repo ports.EmployeeRepository) *employeeService {
	return &employeeService{
		repo: repo,
	}
}

func (s *employeeService) GetEmployee(ctx context.Context, req *entity.GetEmployeeRequest) (entity.GetEmployeeResponse, error) {
	res, err := s.repo.GetEmployee(ctx, req)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *employeeService) UpdateEmployee(ctx context.Context, req *entity.UpdateEmployeeRequest) error {
	err := s.repo.UpdateEmployee(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
