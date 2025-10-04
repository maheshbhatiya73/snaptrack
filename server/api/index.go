package api

import (
	"snaptrack/api/routes"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
    // Auth routes
    authGroup := app.Group("/api/auth")
    RegisterAuthRoutes(authGroup)

    // Other routes
    app.Get("/api/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"message": "Hello from API"})
    })
    app.Get("/api/dashboard/recent-activity", routes.RecentActivity)
    routes.RegisterServerRoutes(app)
    app.Post("/api/local/validate-path", routes.ValidateLocalPath)
    routes.RegisterBackupRoutes(app)
    routes.MonitorRoutes(app)
    routes.RegisterDashboardRoutes(app)
}


func RegisterWebSocketRoutes(app *fiber.App) {
	routes.RegisterWebSocketRoutes(app)
}
