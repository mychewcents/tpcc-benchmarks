package services

import (
	"database/sql"
	"log"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
)

// PaymentService interface to the Payment transaction
type PaymentService interface {
	ProcessTransaction(req *models.Payment) (*models.PaymentOutput, error)
}

type paymentServiceImpl struct {
	db *sql.DB
}

// CreateNewPaymentService creates new payment service
func CreateNewPaymentService(db *sql.DB) PaymentService {
	return &paymentServiceImpl{db: db}
}

func (ps *paymentServiceImpl) ProcessTransaction(req *models.Payment) (result *models.PaymentOutput, err error) {

	log.Printf("Starting the Payment Transaction for: w=%d d=%d c=%d p=%f", req.WarehouseID, req.DistrictID, req.CustomerID, req.Amount)

	log.Printf("Completed the Payment Transaction for: w=%d d=%d c=%d p=%f", req.WarehouseID, req.DistrictID, req.CustomerID, req.Amount)
	return result, nil
}
