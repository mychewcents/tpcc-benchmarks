package performance

import (
	"fmt"
	"log"
	"os"
	"sort"
)

// RecordPerformanceMetrics records and stores the metrics in a CSV file
func RecordPerformanceMetrics(experiment, client int, latencies []float64, path string) error {

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

	outputStr := fmt.Sprintf("%d,%d,%d,%f,%f,%f,%f,%f,%f", experiment, client, processedTxs, elapsedTime/1000, throughput, avgLatency, medianLatency, p95, p99)

	log.Println(outputStr)
	csvFile, err := os.Create(fmt.Sprintf("%s/%d_%d.csv", path, experiment, client))
	if err != nil {
		return err
	}
	defer csvFile.Close()

	if _, err := csvFile.WriteString(outputStr); err != nil {
		return err
	}
	return nil
}
