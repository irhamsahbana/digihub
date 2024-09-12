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

	mrs.Get("/processes", h.getMRSs)
	mrs.Patch("/processes/:id", h.renewOffer)
	mrs.Delete("/processes/:id", h.deleteFollowUp)
}

func (h *mrsHandler) getMRSs(c *fiber.Ctx) error {
	var (
		req = new(entity.GetMRSsRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = m.GetLocals(c)
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetMRSs - Failed to parse request query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.UserId = l.GetUserId()
	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetMRSs - Invalid request payload")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.GetMRSs(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(resp, ""))
}

func (h *mrsHandler) renewOffer(c *fiber.Ctx) error {
	var (
		req = new(entity.RenewWACRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = m.GetLocals(c)
	)

	if err := c.BodyParser(&req.VehicleConditionIds); err != nil {
		log.Warn().Err(err).Msg("handler::RenewWAC - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.UserId = l.GetUserId()
	req.WacId = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::RenewWAC - Invalid request payload")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := h.service.RenewWAC(ctx, req); err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(nil, "Berhasil memperbarui Penawaran"))
}

func (h *mrsHandler) deleteFollowUp(c *fiber.Ctx) error {
	var (
		req = new(entity.DeleteFollowUpRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = m.GetLocals(c)
	)

	req.UserId = l.GetUserId()
	req.WacId = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::DeleteFollowUp - Invalid request payload")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := h.service.DeleteFollowUp(ctx, req); err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(nil, "Berhasil menghapus Follow Up"))
}
