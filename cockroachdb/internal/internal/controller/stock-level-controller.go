package controller

import (
	"bufio"
	"database/sql"
	"log"
	"strconv"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/handler"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/services"
)

// NewOrderControllerImpl provides the interface to call the service
type stockLevelControllerImpl struct {
	s services.StockLevelService
}

// CreateStockLevelController get the new controller to execute the New Order Transaction
func CreateStockLevelController(db *sql.DB) handler.NewTransactionController {
	return &stockLevelControllerImpl{
		s: services.CreateStockLevelService(db),
	}
}

// HandleTransaction performs the transaction and returns the execution result in Boolean
func (slc *stockLevelControllerImpl) HandleTransaction(scanner *bufio.Scanner, args []string) bool {
	wID, _ := strconv.Atoi(args[0])
	dID, _ := strconv.Atoi(args[1])
	threshold, _ := strconv.Atoi(args[2])
	lastNOrders, _ := strconv.Atoi(args[3])

	sl := &models.StockLevel{
		WarehouseID: wID,
		DistrictID:  dID,
		Threshold:   threshold,
		LastNOrders: lastNOrders,
	}

	_, err := slc.s.ProcessTransaction(sl)
	if err != nil {
		log.Printf("error found in the stock level transaction. Err: %v", err)
		return false
	}

	return true
}
