package handler

import (
	"codebase-app/internal/adapter"
	m "codebase-app/internal/middleware"
	"codebase-app/internal/module/common/entity"
	"codebase-app/internal/module/common/ports"
	"codebase-app/internal/module/common/repository"
	"codebase-app/internal/module/common/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type commonHandler struct {
	service ports.CommonService
}

func NewCommonHandler() *commonHandler {
	var (
		handler = new(commonHandler)
		repo    = repository.NewCommonRepository()
		service = service.NewCommonService(repo)
	)

	handler.service = service
	return handler
}

func (h *commonHandler) Register(router fiber.Router) {
	master := router.Group("/masters", m.AuthBearer)

	master.Get("/areas", h.GetAreas)
	master.Get("/potencies", h.GetPotencies)
	master.Get("/vehicle-types", h.GetVehicleTypes)
	master.Get("/employees", h.GetEmployees)
	master.Get("/branches", h.GetBranches)
}

func (h *commonHandler) GetAreas(c *fiber.Ctx) error {
	result, err := h.service.GetAreas(c.Context())
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(result, ""))
}

func (h *commonHandler) GetPotencies(c *fiber.Ctx) error {
	result, err := h.service.GetPotencies(c.Context())
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(result, ""))
}

func (h *commonHandler) GetVehicleTypes(c *fiber.Ctx) error {
	result, err := h.service.GetVehicleTypes(c.Context())
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(result, ""))
}

func (h *commonHandler) GetEmployees(c *fiber.Ctx) error {
	var (
		req = new(entity.GetEmployeesRequest)
		ctx = c.Context()
		l   = m.GetLocals(c)
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Error().Err(err).Msg("handler::GetEmployees - Failed to parse request")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()
	req.UserId = l.UserId

	if err := v.Validate(req); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::GetEmployees - Invalid request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	result, err := h.service.GetEmployees(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(result, ""))
}

func (h *commonHandler) GetBranches(c *fiber.Ctx) error {
	var (
		req = new(entity.GetBranchesRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Error().Err(err).Msg("handler::GetBranches - Failed to parse request")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::GetBranches - Invalid request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	result, err := h.service.GetBranches(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(result, ""))
}
