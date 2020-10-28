package main

import (
	"bufio"
	"database/sql"
	"flag"
	"os"
	"strings"
	"time"
	"sort"
	"fmt"

	caller "github.com/mychewcents/ddbms-project/cockroachdb/internal"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/cdbconn"
)

var db *sql.DB

var (
	experiment  = flag.Int("exp", 0, "Experiment Number")
	client      = flag.Int("client", 0, "Client Number")
	connPtr     = flag.String("host", "localhost", "URL / IP of the DB Server")
	portPtr     = flag.Int("port", 27000, "Port to contact the server's CDB Service")
	dbPtr       = flag.String("database", "defaultdb", "Database to connect")
	usernamePtr = flag.String("username", "root", "Username to connect with")
)

func init() {
	var err error
	flag.Parse()

	if *experiment == 0 || *client == 0 {
		panic("Provide Experiment and Client number to proceed")
	}
	db, err = cdbconn.CreateConnection(*connPtr, *portPtr, *dbPtr, *usernamePtr)
	if err != nil {
		panic(err)
	}
}

func outputStats(latencies []float64) {
	sort.Float64s(latencies)

	elapsedTime := 0.
	processedTxs := len(latencies)
	for _, latency := range latencies {
		elapsedTime += latency
	}
	
	throughput := float64(processedTxs) * 1000 / elapsedTime
	avgLatency := elapsedTime / float64(processedTxs)
	var medianLatency float64
	if processedTxs % 2 > 0 {
		medianLatency = latencies[(processedTxs - 1) / 2]
	} else {
		medianLatency = (latencies[processedTxs / 2] + latencies[processedTxs / 2 - 1]) / 2
	}
	p99 := latencies[99 * processedTxs / 100]
	p95 := latencies[95 * processedTxs / 100]

	outputStr := "Total number of transactions processed: %d\n"
	outputStr += "Total elapsed time: %fs\n"
	outputStr += "Throughput: %f\n"
	outputStr += "Average Latency: %fms\n"
	outputStr += "Median Latency(ms): %fms\n"
	outputStr += "p99 Latency(ms): %fms\n"
	outputStr += "p95 Latency(ms): %fms"

	fmt.Println(
		fmt.Sprintf(outputStr,
			processedTxs,
			elapsedTime / 1000,
			throughput,
			avgLatency,
			medianLatency,
			p99,
			p95,
		),
	)
}

func main() {
	// fmt.Println(*experiment, *client)
	var txArgs []string

	scanner := bufio.NewScanner(os.Stdin)
	var latencies []float64
	for scanner.Scan() {
		txArgs = strings.Split(scanner.Text(), ",")

		start := time.Now()
		status := caller.ProcessRequest(db, scanner, txArgs)
		if status == true {
			latencies = append(latencies, float64(time.Since(start)) / float64(time.Millisecond))
		}
	}
	outputStats(latencies)
}
