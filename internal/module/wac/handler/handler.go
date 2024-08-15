package handler

import (
	"codebase-app/internal/adapter"
	integstorage "codebase-app/internal/integration/localstorage"
	m "codebase-app/internal/middleware"
	"codebase-app/internal/module/wac/entity"
	"codebase-app/internal/module/wac/ports"
	"codebase-app/internal/module/wac/repository"
	"codebase-app/internal/module/wac/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type wachHandler struct {
	service ports.WACService
}

func NewWacHandler(storage integstorage.LocalStorageContract) *wachHandler {
	handler := &wachHandler{}

	repo := repository.NewWACRepository()
	handler.service = service.NewWACService(repo, storage)

	return handler
}

func (h *wachHandler) Register(router fiber.Router) {
	wac := router.Group("/wac", m.AuthBearer)

	wac.Post("/documents",
		m.AuthRole([]string{"service_advisor"}),
		h.createWAC,
	)
	wac.Patch(
		"/documents/:id/offerings",
		m.AuthRole([]string{"service_advisor"}),
		h.OfferWAC,
	)
	wac.Patch(
		"/documents/:id/revenues",
		m.AuthRole([]string{"service_advisor"}),
		h.AddRevenue,
	)

	wac.Get("/documents", h.getWACs)
	wac.Get("/documents/:id", h.getWAC)
}

func (h *wachHandler) createWAC(c *fiber.Ctx) error {
	var (
		req   = new(entity.CreateWACRequest)
		ctx   = c.Context()
		v     = adapter.Adapters.Validator
		local = m.Locals{}
		l     = local.GetLocals(c)
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createWAC - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.UserId = l.GetUserId()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::createWAC - Invalid input")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.CreateWAC(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(resp, ""))
}

func (h *wachHandler) getWACs(c *fiber.Ctx) error {
	var (
		req   = new(entity.GetWACsRequest)
		ctx   = c.Context()
		v     = adapter.Adapters.Validator
		local = m.Locals{}
		l     = local.GetLocals(c)
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::getWAC - Failed to parse request query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()
	req.UserId = l.GetUserId()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::getWAC - Invalid input")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.GetWACs(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *wachHandler) getWAC(c *fiber.Ctx) error {
	var (
		req = new(entity.GetWACRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		Id  = c.Params("id")
	)

	req.Id = Id

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::getWAC - Invalid input")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.GetWAC(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *wachHandler) OfferWAC(c *fiber.Ctx) error {
	var (
		req   = new(entity.OfferWACRequest)
		ctx   = c.Context()
		v     = adapter.Adapters.Validator
		local = m.Locals{}
		l     = local.GetLocals(c)
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::OfferWAC - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.Id = c.Params("id")
	req.UserId = l.GetUserId()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::OfferWAC - Invalid input")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.OfferWAC(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *wachHandler) AddRevenue(c *fiber.Ctx) error {
	var (
		req   = new(entity.AddWACRevenueRequest)
		ctx   = c.Context()
		v     = adapter.Adapters.Validator
		local = m.Locals{}
		l     = local.GetLocals(c)
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::AddRevenue - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.Id = c.Params("id")
	req.UserId = l.GetUserId()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::AddRevenue - Invalid input")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.AddRevenue(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
