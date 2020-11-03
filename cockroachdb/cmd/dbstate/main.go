package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/cdbconn"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/logging"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/statistics/dbstate"
)

var db *sql.DB

var (
	experiment     = flag.Int("exp", 0, "Experiment Number")
	configFilePath = flag.String("config", "", "Path of the DB Server configuration")
	nodeID         = flag.Int("node", 0, "Node ID of the server")
	env            = flag.String("env", "dev", "Provide an env: \"dev\" or \"prod\"")
)

func init() {
	flag.Parse()

	if *experiment == 0 {
		panic("Provide Experiment number to proceed")
	}
	if len(*configFilePath) == 0 {
		panic("config file path cannot be empty. use -config flag")
	}
	if *nodeID < 1 || *nodeID > 5 {
		panic("node id should be between 1 and 5")
	}

	if err := logging.SetupLogOutput("dbstate", "logs", *experiment); err != nil {
		panic(err)
	}
}

func main() {
	c := config.GetConfig(*configFilePath, *nodeID)
	db, err := cdbconn.CreateConnection(c.HostNode)
	if err != nil {
		panic(err)
	}

	if err := dbstate.RecordDBState(db, *experiment, "results/dbstate"); err != nil {
		log.Fatalf("Error in recording the DB State. Err: %v", err)
		fmt.Println(err)
	}
}
