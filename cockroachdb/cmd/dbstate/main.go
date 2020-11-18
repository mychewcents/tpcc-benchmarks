package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	caller "github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/common/config"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/common/logging"
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

	if err := caller.RecordDBState(c, *experiment, "results/dbstate"); err != nil {
		log.Fatalf("error o in recording the DB State. Err: %v", err)
		fmt.Println(err)
	}
}
