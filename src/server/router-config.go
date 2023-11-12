package server

import (
	"github.com/shaileshhb/equisplit/src/controllers"
	"github.com/shaileshhb/equisplit/src/routes/api"
)

func (ser *Server) CreateRouterInstance() {
	ser.InitializeRouter()

	userserv := controllers.NewUserController(ser.DB)
	usercon := api.NewUserRouter(userserv)
	ser.RegisterRoutes([]Controller{usercon})
}
