package server

import (
	"github.com/shaileshhb/equisplit/src/controllers"
	"github.com/shaileshhb/equisplit/src/routes/api"
)

func (ser *Server) CreateRouterInstance() {
	ser.InitializeRouter()

	userserv := controllers.NewuserController(ser.DB)
	usercon := api.NewUserRouter(userserv)
	ser.RegisterRoutes([]Controller{usercon})
}
