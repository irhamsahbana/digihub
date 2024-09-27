package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/app_log/entity"
	"codebase-app/internal/module/app_log/ports"
	"codebase-app/internal/module/app_log/repository"
	"codebase-app/internal/module/app_log/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type appLog struct {
	service ports.AppLogService
}

func NewAppLogHandler() *appLog {
	var (
		repo    = repository.NewAppLogRepository()
		service = service.NewAppLogService(repo)
		handler = new(appLog)
	)
	handler.service = service

	return handler
}

func (h *appLog) Register(router fiber.Router) {
	logging := router.Group("/logs")

	logging.Post("/", h.createLog)
}

func (h *appLog) createLog(c *fiber.Ctx) error {
	var (
		req = new(entity.CreateLogRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Error().Err(err).Msg("handler::createLog - An error occurred")
		return c.JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::createLog - An error occurred")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := h.service.CreateLog(ctx, req); err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}
	return c.JSON(response.Success(nil, ""))
}
