package api

import (
	"snaptrack/api/routes"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {


    app.Get("/api/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Hello from API",
        })
    })
    routes.RegisterServerRoutes(app)
    // Local path validation endpoint (no server id)
    app.Post("/api/local/validate-path", routes.ValidateLocalPath)
    routes.RegisterBackupRoutes(app)
}

func RegisterWebSocketRoutes(app *fiber.App) {
    routes.RegisterWebSocketRoutes(app)
}
