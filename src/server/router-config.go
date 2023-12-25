package server

import (
	"github.com/shaileshhb/equisplit/src/controllers"
	"github.com/shaileshhb/equisplit/src/routes/api"
)

func (ser *Server) CreateRouterInstance() {
	ser.InitializeRouter()

	usercon := controllers.NewUserController(ser.DB, ser.RDB)
	userapi := api.NewUserRouter(usercon, ser.Auth, ser.Log)

	groupcon := controllers.NewGroupController(ser.DB)
	groupapi := api.NewGroupRouter(groupcon, ser.Auth, ser.Log)

	usergroupcon := controllers.NewUserGroupController(ser.DB)
	usergroupapi := api.NewUserGroupRouter(usergroupcon, ser.Auth, ser.Log)

	transactioncon := controllers.NewGroupTransactionController(ser.DB)
	transactionapi := api.NewGroupTransactionRouter(transactioncon, ser.Auth, ser.Log)

	invitationcon := controllers.NewUserInvitationController(ser.DB)
	invitationapi := api.NewUserInvitationRouter(invitationcon, ser.Auth, ser.Log)

	ser.RegisterRoutes([]Controller{userapi, groupapi, usergroupapi, transactionapi, invitationapi})
}
