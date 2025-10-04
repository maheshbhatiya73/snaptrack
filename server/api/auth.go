package api

import (
	"github.com/gofiber/fiber/v2"
	"snaptrack/auth"
)

func RegisterAuthRoutes(router fiber.Router) {
	router.Post("/login", func(c *fiber.Ctx) error {
		type loginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var body loginRequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid request payload",
			})
		}

		if err := auth.PAMAuthenticate(body.Username, body.Password); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Authentication failed",
			})
		}

		if !auth.IsSuperUser(body.Username) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "User does not have sudo permissions",
			})
		}

		token, err := auth.GenerateJWT(body.Username)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Token generation failed",
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Login successful",
			"token":   token,
		})
	})
}
