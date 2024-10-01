package service

import (
	"codebase-app/internal/module/client/entity"
	"codebase-app/internal/module/client/ports"
	"context"
)

var _ ports.ClientService = &clientService{}

type clientService struct {
	repo ports.ClientRepository
}

func NewClientService(repo ports.ClientRepository) *clientService {
	return &clientService{
		repo: repo,
	}
}

func (s *clientService) GetClients(ctx context.Context, req *entity.GetClientsRequest) (entity.GetClientsResponse, error) {
	return s.repo.GetClients(ctx, req)
}
