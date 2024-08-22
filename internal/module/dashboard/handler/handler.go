package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/middleware"
	"codebase-app/internal/module/dashboard/entity"
	"codebase-app/internal/module/dashboard/ports"
	"codebase-app/internal/module/dashboard/repository"
	"codebase-app/internal/module/dashboard/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type dashboardHandler struct {
	service ports.DashboardService
}

func NewDashboardHandler() *dashboardHandler {
	handler := new(dashboardHandler)
	repo := repository.NewDashboardRepository()
	service := service.NewDashboardService(repo)

	handler.service = service
	return handler
}

func (h *dashboardHandler) Register(router fiber.Router) {
	dashboard := router.Group("/dashboard", middleware.AuthBearer)

	dashboard.Get("/lead-trends", h.GetLeadsTrends)
	dashboard.Get("/wac-summaries", h.GetWACSummaries)
}

func (h *dashboardHandler) GetLeadsTrends(c *fiber.Ctx) error {
	var (
		req = new(entity.LeadTrendsRequest)
		ctx = c.Context()
		l   = middleware.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	res, err := h.service.GetLeadsTrends(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(res, ""))
}

func (h *dashboardHandler) GetWACSummaries(c *fiber.Ctx) error {
	var (
		req = new(entity.WACSummaryRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetWACSummaries - failed to parse query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error()))
	}

	req.UserId = l.GetUserId()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetWACSummaries - failed to validate request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.GetWACSummary(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(res, ""))
}
