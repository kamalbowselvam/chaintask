package main

import (
	//"database/sql"
	"database/sql"
	"fmt"

	"github.com/kamalbowselvam/chaintask/internal/core/service"
	"github.com/kamalbowselvam/chaintask/internal/handlers/rest"
	"github.com/kamalbowselvam/chaintask/internal/repositories"
	"github.com/kamalbowselvam/chaintask/server"
	_ "github.com/lib/pq"
)

func main() {

	fmt.Println("Hello World Task")

	var db *sql.DB
	var err error

	connstr := "postgresql://root:secret@localhost:5433/chain_task?sslmode=disable"

	db, err = sql.Open("postgres",connstr)
	if err != nil {
			panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	taskRepository := repositories.NewPersistenceStorage(db)
	taskService := service.NewTaskService(taskRepository)
	taskHandler := rest.NewHttpHandler(taskService)

	server := server.NewServer(taskHandler)
	server.Start(":8080")
}
