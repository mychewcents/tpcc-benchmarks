package controller

import (
	"bufio"
	"database/sql"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/handler"
)

type orderStatusControllerImpl struct {
}

// CreateOrderStatusController creates new controller for order status
func CreateOrderStatusController(db *sql.DB) handler.NewTransactionController {
	return &orderStatusControllerImpl{}
}

func (osc *orderStatusControllerImpl) HandleTransaction(scanner *bufio.Scanner, args []string) bool {

	return true
}
