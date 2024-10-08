package service

import (
	"codebase-app/internal/module/user/entity"
	"codebase-app/internal/module/user/ports"
	"context"
)

var _ ports.UserService = &userService{}

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *userService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Login(ctx context.Context, req *entity.LoginRequest) (entity.LoginResponse, error) {
	return s.repo.Login(ctx, req)
}

func (s *userService) GetProfile(ctx context.Context, req *entity.GetProfileRequest) (entity.GetProfileResponse, error) {
	return s.repo.GetProfile(ctx, req)
}
