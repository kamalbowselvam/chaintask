package main

import (
	//"database/sql"
	"database/sql"

	"log"

	pgadapter "github.com/casbin/casbin-pg-adapter"
	"github.com/kamalbowselvam/chaintask/api"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/server"
	"github.com/kamalbowselvam/chaintask/service"
	"github.com/kamalbowselvam/chaintask/token"
	"github.com/kamalbowselvam/chaintask/util"
	_ "github.com/lib/pq"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	var dbconn *sql.DB
	var err error

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can't load the configuration file")
	}

	dbconn, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic(err)
	}
	adapter, err := pgadapter.NewAdapter(config.DBSource)
	if err != nil {
		panic(err)
	}

	if err = dbconn.Ping(); err != nil {
		panic(err)
	}

	taskRepository := db.NewPersistenceStorage(dbconn)
	taskService := service.NewTaskService(taskRepository)
	tokenMaker, _ := token.NewPasetoMaker(config.TokenSymmetricKey)
	taskHandler := api.NewHttpHandler(taskService, tokenMaker, config)

	server := server.NewServer(taskHandler, adapter)
	server.Start(":8080")
}
