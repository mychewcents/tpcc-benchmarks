package models

// PerformanceMetrics stores the calculated the performance metrics
type PerformanceMetrics struct {
	ProcessedTxs   int
	ElapsedTime    float64
	Throughput     float64
	AverageLatency float64
	MedianLatency  float64
	P95Latency     float64
	P99Latency     float64
}
