package handler

import (
	"codebase-app/internal/middleware"
	"codebase-app/internal/module/dashboard/entity"
	"codebase-app/internal/module/dashboard/ports"
	"codebase-app/internal/module/dashboard/repository"
	"codebase-app/internal/module/dashboard/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
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
}

func (h *dashboardHandler) GetLeadsTrends(c *fiber.Ctx) error {
	var (
		req   = new(entity.LeadTrendsRequest)
		ctx   = c.Context()
		local = middleware.Locals{}
		l     = local.GetLocals(c)
	)

	req.UserId = l.GetUserId()

	res, err := h.service.GetLeadsTrends(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.JSON(response.Success(res, ""))
}
