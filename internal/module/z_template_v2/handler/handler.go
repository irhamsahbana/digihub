package handler

import "github.com/gofiber/fiber/v2"

type xxxHandler struct {
}

func NewXXXHandler() *xxxHandler {
	var handler = new(xxxHandler)

	return handler
}

func (h *xxxHandler) Register(router fiber.Router) {

}
