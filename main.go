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
	app := fiber.New()
	app.Use(cors.New())

	app.Static("/static", config.Config("STATIC_ROOT"))
	app.Static("/media", config.Config("MEDIA_ROOT"))

	database.ConnectDB()

	router.SetupRoutes(app)

	port := config.Config("SERVER_PORT")
	log.Fatal(app.Listen(":" + port))
}
