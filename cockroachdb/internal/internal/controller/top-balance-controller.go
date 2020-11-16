package controller

import (
	"bufio"
	"database/sql"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/handler"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/services"
)

type topBalanceControllerImpl struct {
	s services.TopBalanceService
}

// CreateTopBalanceController creates the new top balance controller
func CreateTopBalanceController(db *sql.DB) handler.NewTransactionController {
	return &topBalanceControllerImpl{
		s: services.CreateTopBalanceService(db),
	}
}

func (tbc *topBalanceControllerImpl) HandleTransaction(scanner *bufio.Scanner, args []string) bool {

	return true
}
