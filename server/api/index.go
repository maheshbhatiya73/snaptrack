package api

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App) {
    app.Get("/api/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Hello from API",
        })
    })
}
