package caller

import (
	"bufio"
	"database/sql"
	"os"
	"strings"
	"time"

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

	// switch transactionArgs[0] {
	// case "N":
	// 	return txRouter.ProcessTransaction(scanner, transactionArgs)
	// case "P":
	// 	return payment.ProcessTransaction(db, nil, transactionArgs[1:])
	// case "D":
	// 	return delivery.ProcessTransaction(db, nil, transactionArgs[1:])
	// case "O":
	// 	return orderstatus.ProcessTransaction(db, nil, transactionArgs[1:])
	// case "S":
	// 	return stocklevel.ProcessTransaction(db, nil, transactionArgs[1:])
	// case "I":
	// 	return popularitem.ProcessTransaction(db, nil, transactionArgs[1:])
	// case "T":
	// 	return topbalance.ProcessTransaction(db, nil)
	// case "R":
	// 	// return relatedcustomer.ProcessTransaction(db, nil, transactionArgs[1:])
	// }
	// return false
}
