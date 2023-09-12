package api

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2/persist"
	"github.com/gin-gonic/gin"
	docs "github.com/kamalbowselvam/chaintask/docs"
	"github.com/kamalbowselvam/chaintask/service"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	adapter    persist.Adapter
	config     util.Config
	service    service.TaskService
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, service service.TaskService, adapter persist.Adapter) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		adapter:    adapter,
		config:     config,
		service:    service,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// FIXME api should be versioned
	docs.SwaggerInfo.BasePath = "/"
	router.POST("/users", server.CreateUser)
	router.POST("/users/login", server.LoginUser)

	authRoutes := router.Group("/").Use(AuthMiddleware(server.tokenMaker))

	authRoutes.GET("/auth", AuthMiddleware(server.tokenMaker),
		func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{})
		},
	)

	authRoutes.GET("/tasks/:id", AuthorizeMiddleware(util.READ, server.adapter), server.GetTask)
	authRoutes.POST("/tasks/", AuthorizeMiddleware(util.WRITE, server.adapter), server.CreateTask)
	authRoutes.DELETE("/tasks/:id", AuthorizeMiddleware(util.DELETE, server.adapter), server.DeleteTask)
	authRoutes.PUT("/tasks/:id", AuthorizeMiddleware(util.UPDATE, server.adapter), server.UpdateTask)
	authRoutes.POST("/projects/", AuthorizeMiddleware(util.WRITE, server.adapter), server.CreateProject)
	server.router = router
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
