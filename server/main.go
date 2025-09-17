package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Static("/", "./web/.output/public")
	app.Get("/api/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello from GoFiber API",
		})
	})

	log.Println("Server running at http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
