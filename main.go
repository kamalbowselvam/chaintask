package main

import (
	"database/sql"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kamalbowselvam/chaintask/api"
	"github.com/kamalbowselvam/chaintask/authorization"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/logger"
	"github.com/kamalbowselvam/chaintask/service"
	"github.com/kamalbowselvam/chaintask/util"
	_ "github.com/lib/pq"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	healthconfig "github.com/tavsec/gin-healthcheck/config"
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
// @name authorization
func main() {

	var dbconn *sql.DB
	var err error

	//aa := zap.NewDevelopmentEncoderConfig()
	//aa.EncodeLevel = zapcore.CapitalColorLevelEncoder

	//logger := zap.New(zapcore.NewCore(
	//	zapcore.NewConsoleEncoder(aa),
	//	zapcore.AddSync(colorable.NewColorableStdout()),
	//	zapcore.DebugLevel,
	//))

	config, err := util.LoadConfig(".")
	if err != nil {
		logger.Fatal("Can't load the configuration file")
	}

	dbconn, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic(err)
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	loaders, err := authorization.Load(config.DBSource, "./config/rbac_model.conf")
	if err != nil {
		panic(err)
	}
	authorizationService, err := authorization.NewCasbinAuthorization(*loaders)
	policyManagementService, _ := authorization.NewCasbinManagement(*loaders)
	if err != nil {
		panic(err)
	}
	// Just assuring that the casbin has at least policies for admin
	policyManagementService.CreateSuperAdminPolicies(util.DEFAULT_SUPER_ADMIN)

	if err = dbconn.Ping(); err != nil {
		panic(err)
	}

	taskRepository := db.NewStore(dbconn)
	taskService := service.NewTaskService(taskRepository, policyManagementService)
	logger.Info("Starting Chain Task SaaS Application")

	server, _ := api.NewServer(config, taskService, authorizationService, policyManagementService)

	sqlCheck := checks.SqlCheck{Sql: dbconn}
	health_config := healthconfig.DefaultConfig()
	health_config.HealthPath = "health"
	healthcheck.New(server.Router, health_config, []checks.Check{sqlCheck})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	server.Start(":" + port)
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)

	if err != nil {
		logger.Info(err.Error())
		logger.Fatal("Cannot create a new migrate instance:")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal("Failed to run migrate up ")
	}

	logger.Info("DB migration Successful")
}
