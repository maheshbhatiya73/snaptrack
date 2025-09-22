package main

import (
	"log"

	"snaptrack/api"
	"snaptrack/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db.Connect()
	db.Init()

	app := fiber.New()
	app.Static("/", "./web/.output/public")
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	api.RegisterRoutes(app)
	api.RegisterAuthRoutes(app)
	api.RegisterWebSocketRoutes(app)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
