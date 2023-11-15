package api

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/shaileshhb/equisplit/src/controllers"
	"github.com/shaileshhb/equisplit/src/models"
	"github.com/shaileshhb/equisplit/src/security"
)

type GroupTransactionRouter interface {
	RegisterRoutes(router fiber.Router)
	add(c *fiber.Ctx) error
	markTransactionPaid(c *fiber.Ctx) error
	delete(c *fiber.Ctx) error
}

type groupTransactionRouter struct {
	con  controllers.GroupTransactionController
	auth security.Authentication
	log  zerolog.Logger
}

// NewGroupTransactionRouter will create new instance of UserGroupHistoryRouter.
func NewGroupTransactionRouter(con controllers.GroupTransactionController, auth security.Authentication, log zerolog.Logger) GroupTransactionRouter {
	return &groupTransactionRouter{
		con:  con,
		auth: auth,
		log:  log,
	}
}

// RegisterRoutes will register routes for user-group router.
func (g *groupTransactionRouter) RegisterRoutes(router fiber.Router) {
	router.Post("/group/:groupId<int>/transaction", g.auth.MandatoryAuthMiddleware, g.add)
	router.Put("/transaction/:transactionId<int>", g.auth.MandatoryAuthMiddleware, g.markTransactionPaid)
	router.Delete("/transaction/:transactionId<id>", g.auth.MandatoryAuthMiddleware, g.delete)
	g.log.Info().Msg("GroupTransaction routes registered")
}

// add will add new transaction for specified group and user.
func (g *groupTransactionRouter) add(c *fiber.Ctx) error {
	transaction := models.GroupTransaction{}

	err := c.BodyParser(&transaction)
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id, err := strconv.Atoi(c.Params("groupId"))
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	transaction.GroupId = uint(id)

	err = g.con.Add(&transaction)
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(nil)
}

// markTransactionPaid will mark the transaction has paid
func (u *groupTransactionRouter) markTransactionPaid(c *fiber.Ctx) error {
	transaction := models.GroupTransaction{}

	transactionId, err := strconv.Atoi(c.Params("transactionId"))
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	transaction.ID = uint(transactionId)

	userInterface := c.Locals("user")
	user := userInterface.(*models.User)

	err = u.con.MarkTransactionPaid(&transaction, user.ID)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusAccepted).JSON(nil)
}

// delete will delete specified user from group
func (u *groupTransactionRouter) delete(c *fiber.Ctx) error {

	transactionId, err := strconv.Atoi(c.Params("transactionId"))
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userInterface := c.Locals("user")
	user := userInterface.(*models.User)

	err = u.con.Delete(uint(user.ID), uint(transactionId))
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusAccepted).JSON(nil)
}
