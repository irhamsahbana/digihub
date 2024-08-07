package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/wac/entity"
	"codebase-app/internal/module/wac/ports"
	"codebase-app/internal/module/wac/repository"
	"codebase-app/internal/module/wac/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type wachHandler struct {
	service ports.WACService
}

func NewWacHandler() *wachHandler {
	handler := &wachHandler{}

	repo := repository.NewWACRepository()
	handler.service = service.NewWACService(repo)

	return handler
}

func (h *wachHandler) Register(router fiber.Router) {
	router.Post("/wac/documents", h.createWAC)
	router.Get("/wac/documents", h.getWAC)
}

func (h *wachHandler) createWAC(c *fiber.Ctx) error {
	var (
		req = &entity.CreateWACRequest{}
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Error().Err(err).Msg("handler::createWAC - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Error().Err(err).Msg("handler::createWAC - Invalid input")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := h.service.CreateWAC(ctx, req); err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(nil, ""))
}

func (h *wachHandler) getWAC(c *fiber.Ctx) error {
	filePath := filepath.Join("storage", "private", "01J4JWY298AMC7S9MTZ5ZBAWDD.png")
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Error().Err(err).Msg("handler::getWAC - File not found")
		return c.Status(fiber.StatusNotFound).SendString("File not found")
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Send(fileBytes)

}
