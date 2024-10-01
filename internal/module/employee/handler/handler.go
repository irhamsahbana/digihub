package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/middleware"
	"codebase-app/internal/module/employee/entity"
	"codebase-app/internal/module/employee/ports"
	"codebase-app/internal/module/employee/repository"
	"codebase-app/internal/module/employee/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type employeeHandler struct {
	service ports.EmployeeService
}

func NewEmployeeHandler() *employeeHandler {
	var (
		repo    = repository.NewEmployeeRepository()
		service = service.NewEmployeeService(repo)
		handler = new(employeeHandler)
	)
	handler.service = service

	return handler
}

func (h *employeeHandler) Register(router fiber.Router) {
	employee := router.Group("/employees", middleware.AuthBearer, middleware.AuthRole([]string{"admin"}))

	employee.Post("/", h.createEmployee)
	employee.Get("/:id", h.getEmployee)
	employee.Patch("/:id", h.updateEmployee)
	employee.Delete("/:id", h.deleteEmployee)
}

func (h *employeeHandler) getEmployee(c *fiber.Ctx) error {
	var (
		req = new(entity.GetEmployeeRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	req.Id = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetEmployee - invalid request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.GetEmployee(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(res, ""))
}

func (h *employeeHandler) updateEmployee(c *fiber.Ctx) error {
	var (
		req = new(entity.UpdateEmployeeRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	req.Id = c.Params("id")
	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::UpdateEmployee - invalid request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::UpdateEmployee - invalid request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err := h.service.UpdateEmployee(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(nil, ""))
}

func (h *employeeHandler) createEmployee(c *fiber.Ctx) error {
	var (
		req = new(entity.CreateEmployeeRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::CreateEmployee - invalid request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::CreateEmployee - invalid request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err := h.service.CreateEmployee(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(nil, ""))
}

func (h *employeeHandler) deleteEmployee(c *fiber.Ctx) error {
	var (
		req = new(entity.DeleteEmployeeRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	req.Id = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::DeleteEmployee - invalid request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err := h.service.DeleteEmployee(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(nil, ""))
}
