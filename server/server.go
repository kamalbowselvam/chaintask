package server

import (

	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/api"
)

// server to serve HTTP request for our booking service
type Server struct {
	taskhandler *api.HttpHandler
	router *gin.Engine
}



func NewServer(handler *api.HttpHandler) *Server {

	server := &Server{taskhandler: handler}
	router := gin.Default()

	tokenMaker := handler.GetTokenMaker()
	
	
	router.POST("/users", server.taskhandler.CreateUser)
	router.POST("/users/login", server.taskhandler.LoginUser)

	authRoutes := router.Group("/").Use(api.AuthMiddleware(*tokenMaker))

	authRoutes.GET("/tasks/:id", server.taskhandler.GetTask)
	authRoutes.POST("/tasks/", server.taskhandler.CreateTask)
	server.router = router
	return server
}

func errorResponse(err error) map[string]interface{} {
	return gin.H{"error":err.Error()}

}

func (server *Server) Start(address string ) error {
		return server.router.Run(address)
}