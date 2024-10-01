package handler

import (
	integstorage "codebase-app/internal/integration/localstorage"

	"codebase-app/internal/adapter"
	"codebase-app/internal/middleware"
	"codebase-app/internal/module/promotion/entity"
	"codebase-app/internal/module/promotion/ports"
	"codebase-app/internal/module/promotion/repository"
	"codebase-app/internal/module/promotion/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type promotionHandler struct {
	service ports.PromotionService
}

func NewPromotionHandler(s integstorage.LocalStorageContract) *promotionHandler {
	var (
		repo    = repository.NewPromotionRepository()
		service = service.NewPromotionService(repo, s)
		handler = new(promotionHandler)
	)
	handler.service = service

	return handler
}

func (h *promotionHandler) Register(router fiber.Router) {
	router.Post("/promotions",
		middleware.AuthBearer,
		middleware.AuthRole([]string{"admin"}),
		h.createPromotion,
	)

	router.Get("/promotions", h.GetPromotions)

	router.Delete("/promotions/:id",
		middleware.AuthBearer,
		middleware.AuthRole([]string{"admin"}),
		h.DeletePromotion,
	)
}

func (h *promotionHandler) createPromotion(c *fiber.Ctx) error {
	var (
		req = new(entity.CreatePromotionRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createPromotion - invalid payload")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		req.RemoveImage()
		log.Warn().Err(err).Any("payload", req).Msg("handler::createPromotion - invalid payload")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := h.service.CreatePromotion(ctx, req); err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(nil, ""))
}

func (h *promotionHandler) GetPromotions(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	promotions, err := h.service.GetPromotions(ctx)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(promotions, ""))
}

func (h *promotionHandler) DeletePromotion(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		req = new(entity.DeletePromotionRequest)
		v   = adapter.Adapters.Validator
	)

	req.Id = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::DeletePromotion - invalid payload")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := h.service.DeletePromotion(ctx, req); err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}
