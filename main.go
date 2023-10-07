package main

import (
	"database/sql"

	"github.com/kamalbowselvam/chaintask/api"
	"github.com/kamalbowselvam/chaintask/authorization"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/service"
	"github.com/kamalbowselvam/chaintask/util"
	_ "github.com/lib/pq"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// @title           Swagger ChainTasks API
// @version         1.0
// @description     This a server helping organizing tasks in a chain.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securitydefinitions.apikey  BearerAuth
// @in header
// @name Authorization
func main() {

	var dbconn *sql.DB
	var err error

	aa := zap.NewDevelopmentEncoderConfig()
	aa.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(aa),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))

	config, err := util.LoadConfig(".")
	if err != nil {
		logger.Fatal("Can't load the configuration file")
	}

	dbconn, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic(err)
	}
	loaders, err := authorization.Load(config.DBSource, "./config/rbac_model.conf", *logger)
	if err != nil {
		panic(err)
	}
	authorizationService, err := authorization.NewCasbinAuthorization(*loaders)
	policyManagementService, _ := authorization.NewCasbinManagement(*loaders)
	if err != nil {
		panic(err)
	}
	// Just assuring that the casbin has at least policies for admin
	policyManagementService.CreateAdminPolicies(util.DEFAULT_ADMIN)

	if err = dbconn.Ping(); err != nil {
		panic(err)
	}

	taskRepository := db.NewStore(dbconn, logger)
	taskService := service.NewTaskService(taskRepository, policyManagementService, logger)
	logger.Info("Starting Chain Task SaaS Application")

	server, _ := api.NewServer(config, taskService, authorizationService, policyManagementService, logger)
	server.Start(":8080")
}
