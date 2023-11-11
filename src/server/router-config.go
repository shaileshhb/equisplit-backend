package server

func (ser *Server) CreateRouterInstance() {
	ser.InitializeRouter()

	ser.RegisterRoutes([]Controller{})
}
