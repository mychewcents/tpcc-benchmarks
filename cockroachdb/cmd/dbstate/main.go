package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/dbstate"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/logging"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/cdbconn"
)

var db *sql.DB

var (
	experiment = flag.Int("exp", 0, "Experiment Number")
	configPath = flag.String("config", "configs/dev/local.json", "Path of the DB Server configuration")
)

func init() {
	var err error
	flag.Parse()

	if *experiment == 0 {
		panic("Provide Experiment and Client number to proceed")
	}
	db, err = cdbconn.CreateConnection(*configPath)
	if err != nil {
		panic(err)
	}

	if err := logging.SetupLogOutput("dbstate", "logs", *experiment); err != nil {
		panic(err)
	}
}

func main() {
	if err := dbstate.RecordDBState(db, *experiment, "results/dbstate"); err != nil {
		log.Fatalf("Error in recording the DB State. Err: %v", err)
		fmt.Println(err)
	}
}
