package handler

import (
	"codebase-app/internal/module/common/ports"
	"codebase-app/internal/module/common/repository"
	"codebase-app/internal/module/common/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
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
	master := router.Group("/masters")

	master.Get("/areas", h.GetAreas)
	master.Get("/potencies", h.GetPotencies)
	master.Get("/vehicle-types", h.GetVehicleTypes)
	master.Get("/employees", h.GetEmployees)
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
	return nil
}
