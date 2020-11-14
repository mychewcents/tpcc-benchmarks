package router

import (
	"bufio"
	"database/sql"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/controller"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/handler"
)

// TransactionRouter stores the handlers
type TransactionRouter struct {
	handlers map[string]handler.NewTransactionController
}

// GetNewRouter creates a new router
func GetNewRouter(db *sql.DB) *TransactionRouter {
	handlers := make(map[string]handler.NewTransactionController)
	handlers["N"] = controller.GetNewNewOrderController(db)
	handlers["S"] = controller.GetNewStockLevelController(db)

	return &TransactionRouter{
		handlers: handlers,
	}
}

// ProcessTransaction processes each transaction upon input
func (tr *TransactionRouter) ProcessTransaction(scanner *bufio.Scanner, args []string) (completed bool) {
	if args[0] = "N" || args[0] = "S" {
		completed = tr.handlers[args[0]].HandleTransaction(scanner, args[1:])
	}

	return completed
}
