package caller

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/common/cdbconn"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/common/config"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/controller"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/router"
)

// ProcessTransactions Calls the required DB function
func ProcessTransactions(c config.Configuration) (latencies []float64) {
	db, err := cdbconn.CreateConnection(c.HostNode)
	if err != nil {
		panic(err)
	}

	var txArgs []string

	txRouter := router.CreateTransactionRouter(db)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		txArgs = strings.Split(scanner.Text(), ",")

		start := time.Now()
		completed := txRouter.ProcessTransaction(scanner, txArgs)
		if completed == true {
			latencies = append(latencies, float64(time.Since(start))/float64(time.Millisecond))
		}
	}

	return latencies
}

// ProcessServerSetupRequest processes the commands for the initialization of the database
func ProcessServerSetupRequest(functionName string, configFilePath, env string, nodeID, experiment int) {
	ssRouter := router.CreateServerSetupRouter(configFilePath, nodeID)

	ssRouter.ProcessServerSetupRequest(functionName, configFilePath, env, nodeID, experiment)
}

// RecordPerformanceMetrics stores the performance metrics for the experiment
func RecordPerformanceMetrics(experiment, client int, latencies []float64, dirPath string) (err error) {
	performanceController := controller.CreatePerformanceController()

	err = performanceController.Record(experiment, client, latencies, dirPath)
	return
}

// RecordDBState records db state
func RecordDBState(c config.Configuration, experiment int, dirPath string) (err error) {
	db, err := cdbconn.CreateConnection(c.HostNode)
	if err != nil {
		panic(err)
	}

	dbController := controller.CreateDatabaseStateController(db)
	err = dbController.CalculateDBState(experiment, dirPath)

	return
}
