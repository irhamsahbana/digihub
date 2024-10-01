package ports

import (
	"codebase-app/internal/module/promotion/entity"
	"context"
)

type PromotionRepository interface {
	CreatePromotion(ctx context.Context, req *entity.CreatePromotionRequest) error
	GetPromotions(ctx context.Context) ([]entity.Promotion, error)
	DeletePromotion(ctx context.Context, req *entity.DeletePromotionRequest) error
}

type PromotionService interface {
	CreatePromotion(ctx context.Context, req *entity.CreatePromotionRequest) error
	GetPromotions(ctx context.Context) ([]entity.Promotion, error)
	DeletePromotion(ctx context.Context, req *entity.DeletePromotionRequest) error
}
