package service

import (
	integstorage "codebase-app/internal/integration/localstorage"

	"codebase-app/internal/module/promotion/entity"
	"codebase-app/internal/module/promotion/ports"
	"context"
)

var _ ports.PromotionService = &promotionService{}

type promotionService struct {
	repo    ports.PromotionRepository
	storage integstorage.LocalStorageContract
}

func NewPromotionService(
	r ports.PromotionRepository,
	s integstorage.LocalStorageContract,
) *promotionService {
	return &promotionService{
		repo:    r,
		storage: s,
	}
}

func (s *promotionService) CreatePromotion(ctx context.Context, req *entity.CreatePromotionRequest) error {
	fullpath, err := s.storage.Save(req.Image, "storage/public/promotions")
	if err != nil {
		return err
	}
	req.Path = fullpath

	err = s.repo.CreatePromotion(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *promotionService) GetPromotions(ctx context.Context) ([]entity.Promotion, error) {
	return s.repo.GetPromotions(ctx)
}

func (s *promotionService) DeletePromotion(ctx context.Context, req *entity.DeletePromotionRequest) error {
	err := s.repo.DeletePromotion(ctx, req)
	if err != nil {
		return err
	}

	err = s.storage.Delete(req.Path)
	if err != nil {
		return err
	}

	return nil
}
