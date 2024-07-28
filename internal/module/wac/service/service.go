package service

import "codebase-app/internal/module/z_template_v2/ports"

type xxxService struct {
	repo ports.XxxRepository
}

func NewXxxService(repo ports.XxxRepository) ports.XxxService {
	return &xxxService{
		repo: repo,
	}
}
