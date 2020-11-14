package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	caller "github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/common/cdbconn"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/common/config"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/common/logging"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/statistics/performance"
)

var db *sql.DB

var (
	experiment     = flag.Int("exp", 0, "Experiment Number")
	client         = flag.Int("client", 0, "Client Number")
	configFilePath = flag.String("config", "", "Path of the DB Server configuration")
	nodeID         = flag.Int("node", 0, "Node ID to run on")
)

func init() {
	flag.Parse()

	if *experiment == 0 || *client == 0 {
		panic("provide Experiment and Client number to proceed")
	}
	if len(*configFilePath) == 0 {
		panic("config file path cannot be empty. use -config option")
	}
	if *nodeID < 1 || *nodeID > 5 {
		panic("node id should be between 1 and 5")
	}

	if err := logging.SetupLogOutput("exp", "logs", *experiment, *client); err != nil {
		panic(err)
	}
}

func main() {
	log.Printf("Starting the experiment: '%d' on node: '%d' ", *experiment, *nodeID)

	c := config.GetConfig(*configFilePath, *nodeID)
	db, err := cdbconn.CreateConnection(c.HostNode)
	if err != nil {
		panic(err)
	}

	latencies := caller.ProcessTransactions(db)
	if len(latencies) == 0 {
		log.Printf("error in the transactions")
		fmt.Println("error occurred in performing transactions. Please check the logs")
		return
	}

	if err := performance.RecordPerformanceMetrics(*experiment, *client, latencies, "results/metrics"); err != nil {
		log.Printf("error in performance recording. Err: %v", err)
		fmt.Println("error occurred in performance recording. Please check the logs")
		return
	}

	log.Printf("Successfully completed the experiment: '%d' on node: '%d' ", *experiment, *nodeID)
}
