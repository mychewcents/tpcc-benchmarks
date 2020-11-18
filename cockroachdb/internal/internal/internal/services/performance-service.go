package services

import (
	"fmt"
	"sort"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
)

// PerformanceService interface to calculate the performance metrics
type PerformanceService interface {
	Calculate(latencies []float64) (*models.PerformanceMetrics, error)
}

type performanceServiceImpl struct {
}

// CreatePerformanceService creates the service to calculate the performance
func CreatePerformanceService() PerformanceService {
	return &performanceServiceImpl{}
}

func (ps *performanceServiceImpl) Calculate(latencies []float64) (result *models.PerformanceMetrics, err error) {

	if len(latencies) == 0 {
		return nil, fmt.Errorf("empty list of latencies passed")
	}

	sort.Float64s(latencies)

	elapsedTime := 0.0
	result.ProcessedTxs = len(latencies)

	for _, latency := range latencies {
		elapsedTime += latency
	}

	result.ElapsedTime = elapsedTime

	result.Throughput = float64(result.ProcessedTxs) * 1000 / elapsedTime
	result.AverageLatency = elapsedTime / float64(result.ProcessedTxs)

	if result.ProcessedTxs%2 == 0 {
		result.MedianLatency = (latencies[result.ProcessedTxs/2] + latencies[result.ProcessedTxs/2-1]) / 2
	} else {
		result.MedianLatency = latencies[(result.ProcessedTxs-1)/2]
	}
	result.P99Latency = latencies[99*result.ProcessedTxs/100]
	result.P95Latency = latencies[95*result.ProcessedTxs/100]

	return
}
