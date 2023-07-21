package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/kamalbowselvam/chaintask/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../")

	if err != nil {
		log.Fatal("Failed to load the config file")
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connet to db: ", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())

}
