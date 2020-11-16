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

type deliveryControllerImpl struct {
	s services.DeliveryService
}

// CreateDeliveryController creates new controller with delivery transaction
func CreateDeliveryController(db *sql.DB) handler.NewTransactionController {
	return &deliveryControllerImpl{
		s: services.CreateDeliveryService(db),
	}
}

func (dc *deliveryControllerImpl) HandleTransaction(scanner *bufio.Scanner, args []string) bool {

	wID, _ := strconv.Atoi(args[0])
	carrierID, _ := strconv.Atoi(args[1])

	d := &models.Delivery{
		WarehouseID: wID,
		CarrierID:   carrierID,
	}

	_, err := dc.s.ProcessTransaction(d)
	if err != nil {
		log.Printf("error occurred in delivery transaction. Err: %v", err)
		return false
	}

	return true
}
