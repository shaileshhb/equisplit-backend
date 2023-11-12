package api

import "github.com/gofiber/fiber/v2"

type GroupRouter interface {
	RegisterRoutes(router *fiber.Router)
}

type groupRouter struct {
}

func NewGroupRouter() UserRouter {
	return &userRouter{
		// service: service,
	}
}

func (g *groupRouter) RegisterRoutes(router *fiber.Router) {
}
