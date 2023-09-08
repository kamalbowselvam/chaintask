package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamalbowselvam/chaintask/authorization"
	docs "github.com/kamalbowselvam/chaintask/docs"
	"github.com/kamalbowselvam/chaintask/service"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	authorize  authorization.AuthorizationService
	policies   authorization.PolicyManagementService
	config     util.Config
	service    service.TaskService
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, service service.TaskService, authorize authorization.AuthorizationService, policies authorization.PolicyManagementService) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		authorize:  authorize,
		policies:   policies,
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

	authRoutes := router.Group("/").Use(AuthMiddleware(server.tokenMaker))

	authRoutes.GET("/auth", AuthMiddleware(server.tokenMaker),
		func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{})
		},
	)
	authorizeMid := AuthorizeMiddleware(server.authorize)
	authRoutes.POST("/users/login", authorizeMid, server.LoginUser)
	authRoutes.GET("/tasks/:id", authorizeMid, server.GetTask)
	authRoutes.POST("/tasks/", authorizeMid, server.CreateTask)
	authRoutes.DELETE("/tasks/:id", authorizeMid, server.DeleteTask)
	authRoutes.PUT("/tasks/:id", authorizeMid, server.UpdateTask)
	authRoutes.POST("/projects/", authorizeMid, server.CreateProject)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
