package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"snaptrack/api"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	app := fiber.New()

	// Serve static files
	app.Static("/", "./web/.output/public")

	// Register APIs
	api.RegisterAuthRoutes(app)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
