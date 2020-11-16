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

type orderStatusControllerImpl struct {
	s services.OrderStatusService
}

// CreateOrderStatusController creates new controller for order status
func CreateOrderStatusController(db *sql.DB) handler.NewTransactionController {
	return &orderStatusControllerImpl{
		s: services.CreateOrderStatusService(db),
	}
}

func (osc *orderStatusControllerImpl) HandleTransaction(scanner *bufio.Scanner, args []string) bool {

	wID, _ := strconv.Atoi(args[0])
	dID, _ := strconv.Atoi(args[1])
	cID, _ := strconv.Atoi(args[2])

	os := &models.OrderStatus{
		WarehouseID: wID,
		DistrictID:  dID,
		CustomerID:  cID,
	}

	_, err := osc.s.ProcessTransaction(os)
	if err != nil {
		log.Printf("error occurred in the order status transaction. Err: %v", err)
		return false
	}

	return true
}
