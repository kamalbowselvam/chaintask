package main

import (
	"database/sql"
	"fmt"
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
	return

}
