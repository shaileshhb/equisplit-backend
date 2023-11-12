package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/shaileshhb/equisplit/src/controllers"
	"github.com/shaileshhb/equisplit/src/models"
)

type GroupRouter interface {
	RegisterRoutes(router *fiber.Router)
}

type groupRouter struct {
	controller controllers.GroupController
}

func NewGroupRouter(controller controllers.GroupController) GroupRouter {
	return &groupRouter{
		controller: controller,
	}
}

func (g *groupRouter) RegisterRoutes(router *fiber.Router) {
	(*router).Post("/:userId/group")
}

func (g *groupRouter) CreateGroup(c *fiber.Ctx) error {
	group := &models.Group{}

	err := c.BodyParser(group)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = g.controller.CreateGroup(group)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User successfully logged in",
	})
}
