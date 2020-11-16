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

type relatedCustomerControllerImpl struct {
	s services.RelatedCustomerService
}

// CreateRelatedCustomerController creates the new controller for the related customer transaction
func CreateRelatedCustomerController(db *sql.DB) handler.NewTransactionController {
	return &relatedCustomerControllerImpl{
		s: services.CreateRelatedCustomerService(db),
	}
}

func (rcc *relatedCustomerControllerImpl) HandleTransaction(scanner *bufio.Scanner, args []string) bool {
	wID, _ := strconv.Atoi(args[0])
	dID, _ := strconv.Atoi(args[1])
	cID, _ := strconv.Atoi(args[2])

	rc := &models.RelatedCustomer{
		WarehouseID: wID,
		DistrictID:  dID,
		CustomerID:  cID,
	}

	_, err := rcc.s.ProcessTransaction(rc)
	if err != nil {
		log.Printf("error occurred in processing the related customer transaction. Err: %v", err)
		return false
	}

	return true
}
