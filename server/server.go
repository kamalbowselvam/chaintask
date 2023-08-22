package server

import (
	"github.com/casbin/casbin/v2/persist"
	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/api"
	"github.com/kamalbowselvam/chaintask/util"
	docs "github.com/kamalbowselvam/chaintask/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// server to serve HTTP request for our booking service
type Server struct {
	taskhandler *api.HttpHandler
	router      *gin.Engine
}

func NewServer(handler *api.HttpHandler, adapter persist.Adapter) *Server {

	server := &Server{taskhandler: handler}
	router := gin.Default()

	tokenMaker := handler.GetTokenMaker()

	// FIXME api should be versioned
	docs.SwaggerInfo.BasePath = "/"
	router.POST("/users", server.taskhandler.CreateUser)
	router.POST("/users/login", server.taskhandler.LoginUser)

	authRoutes := router.Group("/").Use(api.AuthMiddleware(*tokenMaker))

	authRoutes.GET("/tasks/:id", api.AuthorizeMiddleware(util.READ, adapter), server.taskhandler.GetTask)
	authRoutes.POST("/tasks/", api.AuthorizeMiddleware(util.WRITE, adapter), server.taskhandler.CreateTask)
	authRoutes.DELETE("/tasks/:id", api.AuthorizeMiddleware(util.DELETE, adapter), server.taskhandler.DeleteTask)
	authRoutes.PUT("/tasks/:id", api.AuthorizeMiddleware(util.UPDATE, adapter), server.taskhandler.UpdateTask)
	authRoutes.POST("/projects/", api.AuthorizeMiddleware(util.WRITE, adapter), server.taskhandler.CreateProject)
	server.router = router
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
