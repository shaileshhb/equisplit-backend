package security

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Authentication struct {
	log zerolog.Logger
}

func NewAuthentication(log zerolog.Logger) Authentication {
	return Authentication{
		log: log,
	}
}

// MandatoryAuthMiddleware will check that authorization cookie is valid.
func (a *Authentication) MandatoryAuthMiddleware(c *fiber.Ctx) error {
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
func (a *Authentication) OptionalAuthMiddleware(c *fiber.Ctx) error {
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

// HttpLogger will log when an API is called
func (a *Authentication) HttpLogger(c *fiber.Ctx) error {
	a.log.Info().Str("method", c.Method()).
		Str("path", c.OriginalURL()).
		Msg("")

	return c.Next()
}
