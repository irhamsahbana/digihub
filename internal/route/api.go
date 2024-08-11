package route

import (
	integlocalstorage "codebase-app/internal/integration/localstorage"
	m "codebase-app/internal/middleware"
	commonHandler "codebase-app/internal/module/common/handler"
	userHandler "codebase-app/internal/module/user/handler"
	wacHandler "codebase-app/internal/module/wac/handler"

	"codebase-app/pkg/response"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func SetupRoutes(app *fiber.App) {
	storage := integlocalstorage.NewLocalStorageIntegration()

	// add /api prefix to all routes
	api := app.Group("/api")
	api.Get("/storage/private/:filename", m.ValidateSignedURL, storageFile)

	wacHandler.NewWacHandler(storage).Register(api)
	userHandler.NewUserHandler().Register(api)
	commonHandler.NewCommonHandler().Register(api)

	// fallback route
	app.Use(func(c *fiber.Ctx) error {
		var (
			method = c.Method()                       // get the request method
			path   = c.Path()                         // get the request path
			query  = c.Context().QueryArgs().String() // get all query params
			ua     = c.Get("User-Agent")              // get the request user agent
			ip     = c.IP()                           // get the request IP
		)

		log.Trace().
			Str("method", method).
			Str("path", path).
			Str("query", query).
			Str("ua", ua).
			Str("ip", ip).
			Msg("Route not found.")
		return c.Status(fiber.StatusNotFound).JSON(response.Error("Route not found."))
	})
}

func storageFile(c *fiber.Ctx) error {
	var (
		fileName = c.Params("filename")
		filePath = filepath.Join("storage", "private", fileName)
	)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Error().Err(err).Any("url", filePath).Msg("handler::getWAC - File not found")
		return c.Status(fiber.StatusNotFound).JSON(response.Error("File not found"))
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Error().Err(err).Any("url", filePath).Msg("handler::getWAC - Failed to read file")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err.Error()))
	}

	return c.Send(fileBytes)
}
