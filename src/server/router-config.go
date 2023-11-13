package server

import (
	"github.com/shaileshhb/equisplit/src/controllers"
	"github.com/shaileshhb/equisplit/src/routes/api"
)

func (ser *Server) CreateRouterInstance() {
	ser.InitializeRouter()

	usercon := controllers.NewUserController(ser.DB)
	userapi := api.NewUserRouter(usercon)

	groupcon := controllers.NewGroupController(ser.DB)
	groupapi := api.NewGroupRouter(groupcon)

	usergroupcon := controllers.NewUserGroupController(ser.DB)
	usergroupapi := api.NewUserGroupRouter(usergroupcon)
	ser.RegisterRoutes([]Controller{userapi, groupapi, usergroupapi})
}
