package server

import (
	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/internal/handlers/rest"
)

// server to serve HTTP request for our booking service
type Server struct {
	taskhandler *rest.HttpHandler
	router *gin.Engine
}



func NewServer(handler *rest.HttpHandler) *Server {

	server := &Server{taskhandler: handler}
	router := gin.Default()

	router.GET("/tasks/:id", server.taskhandler.GetTask)
	router.POST("/tasks/", server.taskhandler.CreateTask)
	server.router = router
	return server
}

func errorResponse(err error) map[string]interface{} {
	return gin.H{"error":err.Error()}

}

func (server *Server) Start(address string ) error {
		return server.router.Run(address)
}