package caller

import (
	"bufio"
	"database/sql"
	"os"
	"strings"
	"time"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/controller"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/router"
)

// ProcessTransactions Calls the required DB function
func ProcessTransactions(db *sql.DB) (latencies []float64) {
	var txArgs []string

	txRouter := router.GetNewRouter(db)
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

// LoadProcessedTables loads the tables from CSV
func LoadProcessedTables(db *sql.DB) error {
	loadTablesController := controller.CreateLoadProcessedTablesController(db)

	err := loadTablesController.LoadTables()
	if err != nil {
		return err
	}

	return nil
}

// LoadRawTables loads the tables from CSV
func LoadRawTables(db *sql.DB) error {
	loadTablesController := controller.CreateLoadProcessedTablesController(db)

	err := loadTablesController.LoadTables()
	if err != nil {
		return err
	}

	return nil
}
