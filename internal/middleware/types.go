package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Locals struct {
	UserId     string
	Role       string
	IsVerified bool
}

func (l *Locals) GetLocals(c *fiber.Ctx) Locals {
	userId, ok := c.Locals("user_id").(string)
	if ok {
		l.UserId = userId
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get user_id from locals")
	}

	role, ok := c.Locals("role").(string)
	if ok {
		l.Role = role
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get role from locals")
	}

	isVerified, ok := c.Locals("is_verified").(bool)
	if ok {
		l.IsVerified = isVerified
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get is_verified from locals")
	}

	return *l
}

func (l *Locals) GetUserId() string {
	return l.UserId
}

func (l *Locals) GetRole() string {
	return l.Role
}

func (l *Locals) GetIsVerified() bool {
	return l.IsVerified
}
