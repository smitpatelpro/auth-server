package router

import (
	"auth-server/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())
	api.Get("/", handler.HandleRoot)
	api.Post("/signup", handler.Signup)
	api.Post("/login", handler.Login)

}
