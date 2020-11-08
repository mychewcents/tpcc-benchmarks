package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	caller "github.com/mychewcents/ddbms-project/cockroachdb/internal"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/cdbconn"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/logging"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/statistics/performance"
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
	var txArgs []string

	c := config.GetConfig(*configFilePath, *nodeID)
	db, err := cdbconn.CreateConnection(c.HostNode)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	var latencies []float64

	for scanner.Scan() {
		txArgs = strings.Split(scanner.Text(), ",")

		start := time.Now()
		status := caller.ProcessRequest(db, scanner, txArgs)
		if status == true {
			latencies = append(latencies, float64(time.Since(start))/float64(time.Millisecond))
		}
	}

	if err := performance.RecordPerformanceMetrics(*experiment, *client, latencies, "results/metrics"); err != nil {
		log.Printf("error in performance recording. Err: %v", err)
		fmt.Println("error occurred in performance recording. Please check the logs")
		return
	}

	log.Printf("Successfully completed the experiment: '%d' on node: '%d' ", *experiment, *nodeID)
}
