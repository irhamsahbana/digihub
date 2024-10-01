package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/middleware"
	"codebase-app/internal/module/client/entity"
	"codebase-app/internal/module/client/ports"
	"codebase-app/internal/module/client/repository"
	"codebase-app/internal/module/client/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type clientHandler struct {
	service ports.ClientService
}

func NewClientHandler() *clientHandler {
	var (
		repo    = repository.NewClientRepository()
		service = service.NewClientService(repo)
		handler = new(clientHandler)
	)
	handler.service = service

	return handler
}

func (h *clientHandler) Register(router fiber.Router) {
	client := router.Group("/clients", middleware.AuthBearer, middleware.AuthRole([]string{"admin"}))

	client.Get("/", h.GetClients)
}

func (h *clientHandler) GetClients(c *fiber.Ctx) error {
	var (
		req = new(entity.GetClientsRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetClients - failed to parse request")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetClients - invalid request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.GetClients(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(res, ""))
}
