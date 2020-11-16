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

// CreateTransactionRouter creates a new router
func CreateTransactionRouter(db *sql.DB) *TransactionRouter {
	handlers := make(map[string]handler.NewTransactionController)
	handlers["N"] = controller.CreateNewOrderController(db)
	handlers["S"] = controller.CreateStockLevelController(db)
	handlers["P"] = controller.CreatePaymentController(db)
	handlers["I"] = controller.CreatePopularItemController(db)
	handlers["T"] = controller.CreateTopBalanceController(db)
	handlers["D"] = controller.CreateDeliveryController(db)
	handlers["R"] = controller.CreateRelatedCustomerController(db)
	handlers["O"] = controller.CreateOrderStatusController(db)

	return &TransactionRouter{handlers: handlers}
}

// ProcessTransaction processes each transaction upon input
func (tr *TransactionRouter) ProcessTransaction(scanner *bufio.Scanner, args []string) (completed bool) {
	if args[0] == "N" || args[0] == "S" || args[0] == "P" || args[0] == "I" || args[0] == "T" || args[0] == "D" || args[0] == "R" || args[0] == "O" {
		completed = tr.handlers[args[0]].HandleTransaction(scanner, args[1:])
	}

	return completed
}
