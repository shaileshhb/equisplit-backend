package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	router.Post("/group/:groupId<uuid>/transaction", g.auth.MandatoryAuthMiddleware, g.add)
	router.Post("/group/:groupId<uuid>/transactions", g.auth.MandatoryAuthMiddleware, g.addMultiple)
	router.Put("/transaction/:transactionId<uuid>", g.auth.MandatoryAuthMiddleware, g.markTransactionPaid)
	router.Delete("/transaction/:transactionId<uuid>", g.auth.MandatoryAuthMiddleware, g.delete)
	router.Get("/group/:groupId/transactions", g.auth.MandatoryAuthMiddleware, g.getTransactionDetails)
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

	id, err := uuid.Parse(c.Params("groupId"))
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userInterface := c.Locals("user")
	user := userInterface.(*models.User)

	transaction.PayerId = user.Id
	transaction.GroupId = id

	err = transaction.Validate()
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = g.con.Add(&transaction)
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(nil)
}

// addMultiple will add new transaction in specified group.
func (g *groupTransactionRouter) addMultiple(c *fiber.Ctx) error {
	transactions := []models.GroupTransaction{}

	err := c.BodyParser(&transactions)
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	groupId, err := uuid.Parse(c.Params("groupId"))
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userInterface := c.Locals("user")
	user := userInterface.(*models.User)

	for index := range transactions {
		transactions[index].GroupId = groupId
		transactions[index].PayerId = user.Id

		err = transactions[index].Validate()
		if err != nil {
			g.log.Error().Err(err).Msg("")
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	err = g.con.AddMulitple(&transactions)
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(nil)
}

// markTransactionPaid will mark the transaction has paid
func (g *groupTransactionRouter) markTransactionPaid(c *fiber.Ctx) error {
	transaction := models.GroupTransaction{}

	transactionId, err := uuid.Parse(c.Params("transactionId"))
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	transaction.Id = transactionId

	userInterface := c.Locals("user")
	user := userInterface.(*models.User)
	transaction.PayeeId = user.Id

	err = g.con.MarkTransactionPaid(&transaction, user.Id)
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusAccepted).JSON(nil)
}

// delete will delete specified user from group
func (u *groupTransactionRouter) delete(c *fiber.Ctx) error {

	transactionId, err := uuid.Parse(c.Params("transactionId"))
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userInterface := c.Locals("user")
	user := userInterface.(*models.User)

	err = u.con.Delete(user.Id, transactionId)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusAccepted).JSON(nil)
}

// getTransactionDetails will fetch amount to be fetched from all users for specified group
func (u *groupTransactionRouter) getTransactionDetails(c *fiber.Ctx) error {
	var userBalances []models.UserBalance

	groupId, err := uuid.Parse(c.Params("groupId"))
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userInterface := c.Locals("user")
	user := userInterface.(*models.User)

	err = u.con.GetTransactionDetails(&userBalances, user.Id, groupId)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(userBalances)
}
