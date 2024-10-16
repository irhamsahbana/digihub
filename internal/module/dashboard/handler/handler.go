package handler

import (
	"codebase-app/internal/adapter"
	m "codebase-app/internal/middleware"
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
	dashboard := router.Group("/dashboard", m.AuthBearer)

	dashboard.Get("/wac-summaries", h.GetWACSummaries)
	dashboard.Get("/lead-trends", h.GetLeadsTrends)
	dashboard.Get("/admin/summaries", m.AuthRole([]string{"admin"}), h.GetAdminWACSummaries)
	dashboard.Get("/admin/wac-line-chart", m.AuthRole([]string{"admin"}), h.GetWACLineChart)
	dashboard.Get("/admin/activities",
		m.AuthRole([]string{"admin"}),
		h.GetActivities,
	)
}

func (h *dashboardHandler) GetLeadsTrends(c *fiber.Ctx) error {
	var (
		req = new(entity.LeadTrendsRequest)
		ctx = c.Context()
		l   = m.GetLocals(c)
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
		l   = m.GetLocals(c)
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetWACSummaries - failed to parse query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error()))
	}

	req.UserId = l.GetUserId()
	req.UserRole = l.GetRole()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetWACSummaries - failed to validate request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if l.GetRole() == "technician" {
		res, err := h.service.GetWACSummaryTechnician(ctx, req)
		if err != nil {
			code, errs := errmsg.Errors[error](err)
			return c.Status(code).JSON(response.Error(errs))
		}

		return c.JSON(response.Success(res, ""))
	}

	if l.GetRole() == "service_advisor" {
		res, err := h.service.GetWACSummary(ctx, req)
		if err != nil {
			code, errs := errmsg.Errors[error](err)
			return c.Status(code).JSON(response.Error(errs))
		}

		return c.JSON(response.Success(res, ""))
	}

	return c.Status(fiber.StatusForbidden).JSON(response.Error("Forbidden access"))
}

func (h *dashboardHandler) GetWACLineChart(c *fiber.Ctx) error {
	var (
		req = new(entity.GetWACLineChartRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetWACLineChart - failed to parse query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error()))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetWACLineChart - failed to validate request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := req.Validate(); err != nil {
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.GetWACLineChart(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(res, ""))
}

func (h *dashboardHandler) GetActivities(c *fiber.Ctx) error {
	var (
		req = new(entity.GetActivitiesRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetActivities - failed to parse query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error()))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetActivities - failed to validate request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err := req.Validate()
	if err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetActivities - failed to validate request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if req.Export == 1 {
		resp, err := h.service.GetActivitiesExported(ctx, req)
		if err != nil {
			code, errs := errmsg.Errors[error](err)
			return c.Status(code).JSON(response.Error(errs))
		}

		c.Set("Content-Disposition", "attachment; filename=\""+resp.Filename+"\"")
		c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

		return c.SendStream(resp.Buf)
	} else {
		res, err := h.service.GetActivities(ctx, req)
		if err != nil {
			code, errs := errmsg.Errors[error](err)
			return c.Status(code).JSON(response.Error(errs))
		}

		return c.JSON(response.Success(res, ""))
	}
}

func (h *dashboardHandler) GetAdminWACSummaries(c *fiber.Ctx) error {
	var (
		req = new(entity.GetSummaryPerMonthRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetAdminWACSummaries - failed to parse query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error()))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetAdminWACSummaries - failed to validate request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.GetAdminSummary(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(res, ""))
}
