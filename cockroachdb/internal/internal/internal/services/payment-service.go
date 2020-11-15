package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/cockroachdb/cockroach-go/crdb"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dao"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
)

// PaymentService interface to the Payment transaction
type PaymentService interface {
	ProcessTransaction(req *models.Payment) (*models.PaymentOutput, error)
}

type paymentServiceImpl struct {
	db *sql.DB
	w  dao.WarehouseDao
	d  dao.DistrictDao
	c  dao.CustomerDao
}

// CreateNewPaymentService creates new payment service
func CreateNewPaymentService(db *sql.DB) PaymentService {
	return &paymentServiceImpl{
		db: db,
		w:  dao.CreateWarhouseDao(db),
		d:  dao.CreateDistrictDao(db),
		c:  dao.CreateCustomerDao(db),
	}
}

func (ps *paymentServiceImpl) ProcessTransaction(req *models.Payment) (result *models.PaymentOutput, err error) {
	log.Printf("Starting the Payment Transaction for: w=%d d=%d c=%d p=%f", req.WarehouseID, req.DistrictID, req.CustomerID, req.Amount)

	result, err = ps.execute(req)
	if err != nil {
		return nil, fmt.Errorf("error occured while executing the stock level transaction. Err: %v", err)
	}

	log.Printf("Completed the Payment Transaction for: w=%d d=%d c=%d p=%f", req.WarehouseID, req.DistrictID, req.CustomerID, req.Amount)
	return
}

func (ps *paymentServiceImpl) execute(req *models.Payment) (result *models.PaymentOutput, err error) {
	err = crdb.ExecuteTx(context.Background(), ps.db, nil, func(tx *sql.Tx) error {
		result.DistrictAddr, err = ps.d.AddPaymentToDistrict(tx, req.WarehouseID, req.DistrictID, req.Amount)
		if err != nil {
			return err
		}

		result.Customer, err = ps.c.UpdatePaymentDetails(tx, req.WarehouseID, req.DistrictID, req.CustomerID, req.Amount)
		if err != nil {
			return err
		}

		result.WarehouseAddr, err = ps.w.GetAddress(tx, req.WarehouseID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error occurred in updating the tables. Err: %v", err)
	}

	return
}
