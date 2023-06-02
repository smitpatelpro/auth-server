package main

import (
	"auth-server/config"

	"github.com/gofiber/fiber/v2"

	"auth-server/database"
	"auth-server/router"
	"log"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// app := fiber.New()

	// app.Post("/login", handler.Login)

	// // JWT Middleware
	// app.Use(jwtware.New(jwtware.Config{
	// 	SigningKey: []byte("secret"),
	// }))
	// app.Get("/secured", handler.HandleRoot)

	// app.Listen(":3000")

	app := fiber.New()
	app.Use(cors.New())

	app.Static("/static", config.Config("STATIC_ROOT"))
	app.Static("/media", config.Config("MEDIA_ROOT"))
	// utils.CreateMediaDirectories()

	database.ConnectDB()

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
