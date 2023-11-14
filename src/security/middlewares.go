package security

import (
	"github.com/gofiber/fiber/v2"
)

// MandatoryAuthMiddleware will check that authorization cookie is valid.
func MandatoryAuthMiddleware(c *fiber.Ctx) error {
	// authHeader := c.Get("authorization")
	authHeader := c.Cookies("authorization")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	user, err := ValidateJWT(authHeader)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	c.Locals("user", user)
	return c.Next()
}

// OptionalAuthMiddleware will validate authorization cookie only if it exist.
func OptionalAuthMiddleware(c *fiber.Ctx) error {
	// authHeader := c.Cookies("authorization")
	authHeader := c.Get("authorization")

	if authHeader == "" {
		return c.Next()
	}
	user, err := ValidateJWT(authHeader)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	c.Locals("user", user)
	return c.Next()
}
