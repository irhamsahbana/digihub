package ports

import (
	"codebase-app/internal/module/client/entity"
	"context"
)

type ClientRepository interface {
	GetClients(ctx context.Context, req *entity.GetClientsRequest) (entity.GetClientsResponse, error)
}

type ClientService interface {
	GetClients(ctx context.Context, req *entity.GetClientsRequest) (entity.GetClientsResponse, error)
}
