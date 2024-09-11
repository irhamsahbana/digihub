package handler

import (
	"codebase-app/internal/adapter"
	m "codebase-app/internal/middleware"
	"codebase-app/internal/module/mrs/entity"
	"codebase-app/internal/module/mrs/ports"
	"codebase-app/internal/module/mrs/repository"
	"codebase-app/internal/module/mrs/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type mrsHandler struct {
	service ports.MRSService
}

func NewMRSHandler() *mrsHandler {
	var (
		handler = new(mrsHandler)
		repo    = repository.NewMRSRepository()
		service = service.NewMRSService(repo)
	)

	handler.service = service

	return handler
}

func (h *mrsHandler) Register(router fiber.Router) {
	mrs := router.Group("/mrs", m.AuthBearer, m.AuthRole([]string{"technician"}))

	mrs.Get("/processes", h.GetMRSs)
}

func (h *mrsHandler) GetMRSs(c *fiber.Ctx) error {
	var (
		req = new(entity.GetMRSsRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = m.GetLocals(c)
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetMRSs - Failed to parse request query")
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	req.UserId = l.GetUserId()
	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetMRSs - Invalid request payload")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(errs)
	}

	resp, err := h.service.GetMRSs(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(errs)
	}

	return c.JSON(response.Success(resp, ""))
}
