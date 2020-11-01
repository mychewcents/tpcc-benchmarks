package controller

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

type PerformanceMonitorController interface {
	StoreLatency(latency int)
	StorePerformanceMetrics(path string, experimentNo int, clientNo int)
}

type performanceMonitorControllerImpl struct {
	latencies []int
}

func NewPerformanceMonitorController() PerformanceMonitorController {
	return &performanceMonitorControllerImpl{latencies: make([]int, 0)}
}

func (p *performanceMonitorControllerImpl) StoreLatency(latency int) {
	p.latencies = append(p.latencies, latency)
}

func (p *performanceMonitorControllerImpl) StorePerformanceMetrics(path string, experimentNo int, clientNo int) {
	sort.Ints(p.latencies)

	noOfTransactions := len(p.latencies)

	totalLatency := 0
	for _, l := range p.latencies {
		totalLatency += l
	}
	totalLatencyInSeconds := totalLatency / 1000

	throughPut := float64(noOfTransactions) / float64(totalLatencyInSeconds)

	avgLatency := float64(totalLatency) / float64(noOfTransactions)

	var medianTransactionLatency float64
	if len(p.latencies)%2 == 0 {
		medianTransactionLatency = float64(p.latencies[len(p.latencies)/2-1]+p.latencies[len(p.latencies)/2]) / 2.0
	} else {
		medianTransactionLatency = float64(p.latencies[len(p.latencies)/2])
	}

	percentile95TransactionLatency := p.latencies[int(math.Ceil(95.0/100.0*float64(len(p.latencies))))-1]
	percentile99TransactionLatency := p.latencies[int(math.Ceil(99.0/100.0*float64(len(p.latencies))))-1]

	metrics := fmt.Sprintf("%d,%d,%d,%d,%.2f,%.2f,%.2f,%d,%d", experimentNo, clientNo, noOfTransactions, totalLatencyInSeconds,
		throughPut, avgLatency, medianTransactionLatency, percentile95TransactionLatency, percentile99TransactionLatency)

	fileName := fmt.Sprintf("%v/experiment_%v_client_%v.csv", path, experimentNo, clientNo)
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("ERROR saving results, err=%v", err)
	}
	defer file.Close()

	file.WriteString(metrics)
}
