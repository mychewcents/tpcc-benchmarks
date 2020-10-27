package main

import (
	"bufio"
	"database/sql"
	"flag"
	"os"
	"strings"
	"time"
	"sort"

	caller "github.com/mychewcents/ddbms-project/cockroachdb/internal"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/cdbconn"
)

var db *sql.DB

var (
	connPtr     = flag.String("host", "localhost", "URL / IP of the DB Server")
	portPtr     = flag.Int("port", 26257, "Port to contact the server's CDB Service")
	dbPtr       = flag.String("database", "defaultdb", "Database to connect")
	usernamePtr = flag.String("username", "root", "Username to connect with")
)

func init() {
	var err error
	flag.Parse()

	db, err = cdbconn.CreateConnection(*connPtr, *portPtr, *dbPtr, *usernamePtr)
	if err != nil {
		panic(err)
	}
}

func outputStats(latencies []int64) {
	sort.Slice(latencies)
	
	elapsedTime := 0
	processedTxs = len(times)
	for latency := range latencies {
		elapsedTime += latency
	}
	
	throughput := processedTxs / elapsedTime
	avgLatency := elapsedTime / processedTxs
	
	if processedTxs % 2 {
		medianLatency := latencies[(processedTxs - 1) / 2]
	} else {
		medianLatency := (latencies[processedTxs / 2 + 1] + latencies[processedTxs / 2]) / 2
	}
	p99 := latencies[int(.99 * processedTxs)]
	p95 := latencies[int(.95 * processedTxs)]

	outputStr := "Total number of transactions processed: %d\n"
	outputStr += "Total elapsed time(s): %f\n"
	outputStr += "Throughput: %f\n"
	outputStr += "Average Latency(ms): %f\n"
	outputStr += "Median Latency(ms): %f\n"
	outputStr += "p99 Latency(ms): %f\n"
	outputStr += "p95 Latency(ms): %f"

	fmt.Println(outputStr,
		processedTxs,
		elapsedTime / 1000,
		throughput,
		avgLatency,
		medianLatency,
		p99,
		p95
	)
}

func main() {
	var txArgs []string

	scanner := bufio.NewScanner(os.Stdin)
	var processingTimes = []int64
	for scanner.Scan() {
		txArgs = strings.Split(scanner.Text(), ",")

		start := time.Now()
		status := caller.ProcessRequest(db, scanner, txArgs)
		if status == true {
			processingTimes = append(processingTimes, 1. * time.Since(start) / time.Millisecond)
		}
	}

	outputStats(processingTimes)
}
