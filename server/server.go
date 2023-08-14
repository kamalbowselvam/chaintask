package server

import (
	"github.com/casbin/casbin/v2/persist"
	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/internal/handlers/rest"
	"github.com/kamalbowselvam/chaintask/middlewares"
	"github.com/kamalbowselvam/chaintask/util"
)

// server to serve HTTP request for our booking service
type Server struct {
	taskhandler *rest.HttpHandler
	router      *gin.Engine
}

func NewServer(handler *rest.HttpHandler, adapter persist.Adapter) *Server {

	server := &Server{taskhandler: handler}
	router := gin.Default()

	api := router.Group("api/v1")
	api.POST("/users/", server.taskhandler.CreateUser)
	api.POST("/token", server.taskhandler.Login)
	secured := api.Group("secured").Use(middlewares.Auth())
	secured.GET("/tasks/:id", middlewares.Authorize(util.TASK, util.READ, adapter), server.taskhandler.GetTask)
	secured.POST("/tasks/", middlewares.Authorize(util.TASK, util.WRITE, adapter), server.taskhandler.CreateTask)
	server.router = router
	return server
}

func errorResponse(err error) map[string]interface{} {
	return gin.H{"error": err.Error()}

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
