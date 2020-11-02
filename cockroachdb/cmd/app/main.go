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

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/logging"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/performance"

	caller "github.com/mychewcents/ddbms-project/cockroachdb/internal"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/cdbconn"
)

var db *sql.DB

var (
	experiment = flag.Int("exp", 0, "Experiment Number")
	client     = flag.Int("client", 0, "Client Number")
	configPath = flag.String("config", "configs/dev/local.json", "Path of the DB Server configuration")
)

func init() {
	var err error
	flag.Parse()

	if *experiment == 0 || *client == 0 {
		panic("Provide Experiment and Client number to proceed")
	}
	db, err = cdbconn.CreateConnection(*configPath)
	if err != nil {
		panic(err)
	}

	if err := logging.SetupLogOutput("exp", "logs", *experiment, *client); err != nil {
		panic(err)
	}
}

func main() {
	var txArgs []string

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
		log.Fatalf("Error in performance recording. Err: %v", err)
		fmt.Println(err)
	}
}
