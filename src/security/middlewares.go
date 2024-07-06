package security

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Authentication struct {
	// rdb                     *redis.Client
	log                     zerolog.Logger
	authorizationTypeBearer string
}

func NewAuthentication(log zerolog.Logger) Authentication {
	return Authentication{
		// rdb:                     rdb,
		log:                     log,
		authorizationTypeBearer: "bearer",
	}
}

// MandatoryAuthMiddleware will check that authorization cookie is valid.
func (a *Authentication) MandatoryAuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("authorization")

	// authHeader := c.Cookies("authorization")

	if authHeader == "" {
		a.log.Error().Msg("authorization token not specified")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		a.log.Error().Err(errors.New("invalid authorization header provided")).Msg("")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid authorization header provided",
		})
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != a.authorizationTypeBearer {
		a.log.Error().Err(fmt.Errorf("unsupported authorization type %s", authorizationType)).Msg("")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fmt.Sprintf("unsupported authorization type %s", authorizationType),
		})
	}

	user, err := ValidateJWT(fields[1])
	if err != nil {
		a.log.Error().Err(err).Msg("")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
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

	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		a.log.Error().Err(errors.New("invalid authorization header provided")).Msg("")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid authorization header provided",
		})
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != a.authorizationTypeBearer {
		a.log.Error().Err(fmt.Errorf("unsupported authorization type %s", authorizationType)).Msg("")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fmt.Sprintf("unsupported authorization type %s", authorizationType),
		})
	}

	user, err := ValidateJWT(authHeader)
	if err != nil {
		a.log.Error().Err(err).Msg("")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
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
