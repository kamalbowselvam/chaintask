package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	ginprometheus "github.com/zsais/go-gin-prometheus"

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
	Router     *gin.Engine
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("user_role", validRole)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// Prometheus part
	p := ginprometheus.NewPrometheus("chaintask")

	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.Path
		params := [3]string{"companyId", "projectId", "taskId"}
		for _, par := range params {
			for _, p := range c.Params {
				if p.Key == par {
					url = strings.Replace(url, p.Value, par, 1)
					break
				}
			}
		}
		return url
	}

	p.Use(router)

	prometheus_push_url := server.config.PushGateway;
	if(prometheus_push_url != ""){
		p.SetPushGateway(prometheus_push_url, p.MetricsPath, server.config.PushGatewayRefreshDuration)
	}

	// FIXME api should be versioned
	docs.SwaggerInfo.BasePath = "/"
	host := (os.Getenv("HOST_URL"))
	if host != "" {
		docs.SwaggerInfo.Host = host
	}

	globalGroup := router.Group("/").Use(requestLogger())
	globalGroup.POST("/users/login", server.LoginUser)
	globalGroup.POST("/tokens/renew_access", server.renewAccessToken)
	authRoutes := globalGroup.Use(AuthMiddleware(server.tokenMaker))

	authRoutes.GET("/auth", AuthMiddleware(server.tokenMaker),
		func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{})
		},
	)
	authorizeMid := AuthorizeMiddleware(server.authorize)
	authRoutes.POST("/users", authorizeMid, server.CreateUser)
	authRoutes.POST("/company/:companyId/projects/", authorizeMid, server.CreateProject)
	authRoutes.POST("/company/:companyId/projects/:projectId/tasks/", authorizeMid, server.CreateTask)
	authRoutes.POST("/company/:companyId/projects/:projectId/payments/:taskId", authorizeMid, server.PayForATask)
	authRoutes.GET("/company/:companyId/projects/:projectId/tasks/:taskId", authorizeMid, server.GetTask)
	authRoutes.GET("/company/:companyId/projects/:projectId/tasks/", authorizeMid, server.GetTaskListByProject)
	authRoutes.PUT("/company/:companyId/projects/:projectId/tasks/:taskId", authorizeMid, server.UpdateTask)
	authRoutes.DELETE("/company/:companyId/projects/:projectId/tasks/:taskId", authorizeMid, server.DeleteTask)
	authRoutes.POST("/companies/", authorizeMid, server.CreateCompany)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	server.Router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}

func errorResponse(err error) map[string]interface{} {

	return gin.H{"error": err.Error()}

}
