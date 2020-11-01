package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

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
	if processedTxs%2 > 0 {
		medianLatency = latencies[(processedTxs-1)/2]
	} else {
		medianLatency = (latencies[processedTxs/2] + latencies[processedTxs/2-1]) / 2
	}
	p99 := latencies[99*processedTxs/100]
	p95 := latencies[95*processedTxs/100]

	// var outputStr strings.Builder
	// outputStr.WriteString(fmt.Sprintf("Total number of transactions processed: %d\n", processedTxs))
	// outputStr.WriteString(fmt.Sprintf("Total elapsed time: %fs\n", elapsedTime/1000))
	// outputStr.WriteString(fmt.Sprintf("Throughput: %f\n", throughput))
	// outputStr.WriteString(fmt.Sprintf("Average Latency: %fms\n", avgLatency))
	// outputStr.WriteString(fmt.Sprintf("Median Latency(ms): %fms\n", medianLatency))
	// outputStr.WriteString(fmt.Sprintf("p99 Latency(ms): %fms\n", p99))
	// outputStr.WriteString(fmt.Sprintf("p95 Latency(ms): %fms", p95))

	outputStr := fmt.Sprintf("%d,%d,%d,%f,%f,%f,%f,%f,%f", experiment, client, processedTxs, elapsedTime/1000, throughput, avgLatency, medianLatency, p99, p95)

	csvFile, err := os.Create(fmt.Sprintf("metrics/%d_%d.csv", *experiment, *client))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvFile.Close()

	if _, err := csvFile.WriteString(outputStr); err != nil {
		fmt.Println(err)
	}
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
			latencies = append(latencies, float64(time.Since(start))/float64(time.Millisecond))
		}
	}
	outputStats(latencies)
}
