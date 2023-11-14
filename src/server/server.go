package server

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/shaileshhb/equisplit/src/models"
	"gorm.io/gorm"
)

// Controller is implemented by the controllers.
type Controller interface {
	RegisterRoutes(router *fiber.Router)
}

// ModuleConfig needs to be implemented by every module.
type ModuleConfig interface {
	TableMigration(wg *sync.WaitGroup)
}

// Server Struct For Start the equisplit service.
type Server struct {
	Name   string
	DB     *gorm.DB
	App    *fiber.App
	Router *fiber.Router
	WG     *sync.WaitGroup
	// Log    log.Logger
	// Config config.ConfReader
}

func NewServer(name string, db *gorm.DB, wg *sync.WaitGroup) *Server {
	return &Server{
		Name: name,
		DB:   db,
		WG:   wg,
		// Log:  log,
		// Config:         conf,
	}
}

// InitializeRouter Register the route.
func (ser *Server) InitializeRouter() {
	app := fiber.New(fiber.Config{
		AppName: ser.Name,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "Hello world!!",
		})
	})

	apiV1 := app.Group("api/v1")

	ser.App = app
	ser.Router = &apiV1
}

// RegisterRoutes will register the specified routes in controllers.
func (ser *Server) RegisterRoutes(controllers []Controller) {
	for _, controller := range controllers {
		controller.RegisterRoutes(ser.Router)
	}
}

// MigrateTables will do a table table migration for all modules.
func (ser *Server) MigrateTables() {
	lo.Must0(ser.DB.AutoMigrate(&models.User{}))
	lo.Must0(ser.DB.AutoMigrate(&models.Group{}))
	lo.Must0(ser.DB.AutoMigrate(&models.UserGroup{}))
	lo.Must0(ser.DB.AutoMigrate(&models.UserGroupHistory{}))

	config := models.NewModuleConfig(ser.DB)
	config.TableMigration(ser.WG)
	// ser.WG.Add(len(configs))
	// for _, config := range configs {
	// 	config.TableMigration(ser.WG)
	// 	ser.WG.Done()
	// }
	// ser.WG.Wait()
}
