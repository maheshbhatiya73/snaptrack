package api

import (
	"github.com/gofiber/fiber/v2"
	"snaptrack/auth"
)

func RegisterAuthRoutes(app *fiber.App) {
	app.Post("/api/login", func(c *fiber.Ctx) error {
		type loginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var body loginRequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		if err := auth.PAMAuthenticate(body.Username, body.Password); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authentication failed"})
		}

		if !auth.IsSuperUser(body.Username) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not authorized"})
		}

		token, err := auth.GenerateJWT(body.Username)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Token generation failed"})
		}

		return c.JSON(fiber.Map{"token": token})
	})
}
