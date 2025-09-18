package routes

import (
	"snaptrack/auth"
	"snaptrack/db"

	"github.com/gofiber/fiber/v2"
)

func RegisterServerRoutes(app *fiber.App) {
	api := app.Group("/api/servers", auth.RequireJWT())

	api.Get("/", listServers)
	api.Get("/:id", getServer)
	api.Post("/", createServer)
	api.Put("/:id", updateServer)
	api.Delete("/:id", deleteServer)
}

func listServers(c *fiber.Ctx) error {
	var servers []db.Server
	db.DB.Find(&servers)
	return c.JSON(servers)
}

func getServer(c *fiber.Ctx) error {
	id := c.Params("id")
	var server db.Server
	if err := db.DB.First(&server, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Server not found"})
	}
	return c.JSON(server)
}

func createServer(c *fiber.Ctx) error {
	var server db.Server
	if err := c.BodyParser(&server); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := db.DB.Create(&server).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(server)
}

func updateServer(c *fiber.Ctx) error {
	id := c.Params("id")
	var server db.Server
	if err := db.DB.First(&server, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Server not found"})
	}

	var updateData db.Server
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	db.DB.Model(&server).Updates(updateData)
	return c.JSON(server)
}

func deleteServer(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.DB.Delete(&db.Server{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
