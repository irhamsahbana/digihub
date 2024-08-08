package ports

import (
	"codebase-app/internal/module/user/entity"
	"context"
)

type UserRepository interface {
	Login(ctx context.Context, req *entity.LoginRequest) (entity.LoginResponse, error)
}

type UserService interface {
	Login(ctx context.Context, req *entity.LoginRequest) (entity.LoginResponse, error)
}
