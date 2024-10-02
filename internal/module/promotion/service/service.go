package service

import (
	integstorage "codebase-app/internal/integration/localstorage"
	"strings"

	"codebase-app/internal/module/promotion/entity"
	"codebase-app/internal/module/promotion/ports"
	"context"

	"github.com/rs/zerolog/log"
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

func (s *promotionService) UpdatePromotion(ctx context.Context, req *entity.UpdatePromotionRequest) error {
	oldPromotion, err := s.repo.GetPromotionById(ctx, req.Id)
	if err != nil {
		return err
	}

	if req.Image.Present && req.Image.Valid {
		fullpath, err := s.storage.Save(req.Image.Val, "storage/public/promotions")
		if err != nil {
			return err
		}
		req.Path = fullpath
	}

	err = s.repo.UpdatePromotion(ctx, req)
	if err != nil {
		return err
	}

	if req.Image.Present && req.Image.Valid {
		oldPath := strings.ReplaceAll(oldPromotion.Path, "storage/", "./storage/")
		err = s.storage.Delete(oldPath)
		if err != nil {
			log.Warn().Err(err).Str("path", oldPath).Msg("service::UpdatePromotion - failed to delete old image")
		}
	}

	return nil
}
