package main

import (
	//"database/sql"
	"database/sql"

	"log"

	pgadapter "github.com/casbin/casbin-pg-adapter"
	"github.com/kamalbowselvam/chaintask/api"
	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/service"
	"github.com/kamalbowselvam/chaintask/util"
	_ "github.com/lib/pq"
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

	taskRepository := db.NewStore(dbconn)
	taskService := service.NewTaskService(taskRepository)

	server, _ := api.NewServer(config, taskService, adapter)
	server.Start(":8080")
}
