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

// DeliveryService interface to the Delivery Transaction service
type DeliveryService interface {
	ProcessTransaction(req *models.Delivery) (*models.DeliveryOutput, error)
}

type deliveryServiceImpl struct {
	db *sql.DB
	o  dao.OrderDao
	c  dao.CustomerDao
}

// CreateDeliveryService creates the service for Delivery Transaction
func CreateDeliveryService(db *sql.DB) DeliveryService {
	return &deliveryServiceImpl{
		db: db,
		o:  dao.CreateOrderDao(db),
		c:  dao.CreateCustomerDao(db),
	}
}

func (ds *deliveryServiceImpl) ProcessTransaction(req *models.Delivery) (result *models.DeliveryOutput, err error) {
	log.Printf("Starting the Delivery Transaction for: w=%d c=%d", req.WarehouseID, req.CarrierID)

	result, err = ds.execute(req)
	if err != nil {
		return nil, fmt.Errorf("error occurred in processing the delivery transaction. Err: %v", err)
	}

	log.Printf("Completed the Delivery Transaction for: w=%d c=%d", req.WarehouseID, req.CarrierID)
	return
}

func (ds *deliveryServiceImpl) execute(req *models.Delivery) (result *models.DeliveryOutput, err error) {
	orderIDs, err := ds.o.GetUndeliveredOrderIDsPerDistrict(req.WarehouseID)
	if err != nil {
		return nil, err
	}

	err = crdb.ExecuteTx(context.Background(), ds.db, nil, func(tx *sql.Tx) error {
		var totalAmount float64
		var customerID int

		for dID := 1; dID <= 10; dID++ {
			if orderIDs[dID] > 0 {

				if customerID, totalAmount, err = ds.o.DeliverOrder(tx, req.WarehouseID, dID, orderIDs[dID], req.CarrierID); err != nil {
					return fmt.Errorf("error occurred while updating the order details. Err: %v", err)
				}
				if err := ds.c.DeliverOrder(tx, req.WarehouseID, dID, customerID, totalAmount); err != nil {
					return fmt.Errorf("error occurred while updating the customer details. Err: %v", err)
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error occurred while updating the order/customer table. Err: %v", err)
	}

	return
}
