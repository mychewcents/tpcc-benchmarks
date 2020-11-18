package controller

import (
	"fmt"
	"log"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/helper"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/services"
)

// PerformanceController interface to rcord the performances
type PerformanceController interface {
	Record(experiment, client int, latencies []float64, dirPath string) error
}

type performanceControllerImpl struct {
	s services.PerformanceService
}

// CreatePerformanceController creates a new controller to record the performance
func CreatePerformanceController() PerformanceController {
	return &performanceControllerImpl{
		s: services.CreatePerformanceService(),
	}
}

func (pc *performanceControllerImpl) Record(experiment, client int, latencies []float64, dirPath string) (err error) {
	result, err := pc.s.Calculate(latencies)
	if err != nil {
		return err
	}

	outputCSVString := fmt.Sprintf("%d,%d,%d,%f,%f,%f,%f,%f,%f",
		experiment,
		client,
		result.ProcessedTxs,
		result.ElapsedTime,
		result.Throughput,
		result.AverageLatency,
		result.MedianLatency,
		result.P95Latency,
		result.P99Latency,
	)

	finalCSVPath := fmt.Sprintf("%s/%d_%d.csv", dirPath, experiment, client)
	err = helper.WriteFile(outputCSVString, finalCSVPath)
	if err != nil {
		log.Printf("DB State: %s", outputCSVString)
		log.Fatalf("error occurred in writing the csv file. Err: %v", err)
		return err
	}

	return
}
