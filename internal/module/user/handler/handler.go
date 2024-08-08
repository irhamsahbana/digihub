package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/user/entity"
	"codebase-app/internal/module/user/ports"
	"codebase-app/internal/module/user/repository"
	"codebase-app/internal/module/user/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type userHandler struct {
	service ports.UserService
}

func NewUserHandler() *userHandler {
	var (
		handler = new(userHandler)
		repo    = repository.NewUserRepository()
		service = service.NewUserService(repo)
	)

	handler.service = service

	return handler
}

func (h *userHandler) Register(router fiber.Router) {
	auth := router.Group("/authentications/login")

	auth.Post("/", h.login)
}

func (h *userHandler) login(c *fiber.Ctx) error {
	var (
		req = new(entity.LoginRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Error().Err(err).Msg("handler::login - body parser")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Error().Err(err).Msg("handler::login - validation")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.Login(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(res, ""))
}
