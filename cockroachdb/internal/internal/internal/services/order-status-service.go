package services

import (
	"database/sql"
	"log"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
)

// OrderStatusService interface to processing the transactions
type OrderStatusService interface {
	ProcessTransaction(req *models.OrderStatus) (result *models.OrderStatusOutput, err error)
}

type orderStatusServiceImpl struct {
	db *sql.DB
}

// CreateOrderStatusService creates the service for the Order Status Transaction
func CreateOrderStatusService(db *sql.DB) OrderStatusService {
	return &orderStatusServiceImpl{db: db}
}

func (oss *orderStatusServiceImpl) ProcessTransaction(req *models.OrderStatus) (result *models.OrderStatusOutput, err error) {
	log.Printf("Starting the Order Status Transaction for: w=%d d=%d c=%d", req.WarehouseID, req.DistrictID, req.CustomerID)

	result, err = oss.execute(req)
	if err != nil {
		log.Printf("error occured while executing the order status transaction. Err: %v", err)
		return nil, err
	}

	log.Printf("Completed the Order Status Transaction for: w=%d d=%d c=%d", req.WarehouseID, req.DistrictID, req.CustomerID)
	return result, nil
}

func (oss *orderStatusServiceImpl) execute(req *models.OrderStatus) (result *models.OrderStatusOutput, err error) {
	return nil, nil
}
