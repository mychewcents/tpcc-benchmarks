package services

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dao"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
)

// PopularItemService create the interface to process the popular item tx
type PopularItemService interface {
	ProcessTransaction(req *models.PopularItem) (*models.PopularItemOutput, error)
}

type popularItemServiceImpl struct {
	db *sql.DB
	d  dao.DistrictDao
	o  dao.OrderDao
	c  dao.CustomerDao
	i  dao.ItemDao
	ol dao.OrderLineDao
}

// CreateNewPopularItemService creates the new object for the popular item tx
func CreateNewPopularItemService(db *sql.DB) PopularItemService {
	return &popularItemServiceImpl{
		db: db,
		d:  dao.CreateDistrictDao(db),
		o:  dao.CreateOrderDao(db),
		c:  dao.CreateCustomerDao(db),
		i:  dao.CreateItemDao(db),
		ol: dao.CreateOrderLineDao(db),
	}
}

// ProcessTransaction processes the popular item transaction
func (pis *popularItemServiceImpl) ProcessTransaction(req *models.PopularItem) (result *models.PopularItemOutput, err error) {
	log.Printf("Starting the Popular Item Transaction for: w=%d d=%d n=%d", req.WarehouseID, req.DistrictID, req.LastNOrders)

	result, err = pis.execute(req)
	if err != nil {
		return nil, fmt.Errorf("error occurred while executing the popular item transaction. Err: %v", err)
	}

	log.Printf("Completed the Popular Item Transaction for: w=%d d=%d n=%d", req.WarehouseID, req.DistrictID, req.LastNOrders)
	return
}

func (pis *popularItemServiceImpl) execute(req *models.PopularItem) (result *models.PopularItemOutput, err error) {

	result.LastOrderID, err = pis.d.GetLastOrderID(req.WarehouseID, req.DistrictID)
	if err != nil {
		return nil, err
	}

	result.StartOrderID = result.LastOrderID - req.LastNOrders

	orderMap, err := pis.ol.GetMaxQuantityOrderLinesPerOrder(req.WarehouseID, req.DistrictID, result.StartOrderID, result.LastOrderID)
	if err != nil {
		return nil, err
	}

	for key, value := range orderMap {
		orderDetails, err := pis.o.GetDetails(req.WarehouseID, req.DistrictID, key)
		if err != nil {
			return nil, err
		}

		customer, err := pis.c.GetDetails(req.WarehouseID, req.DistrictID, orderDetails.CustomerID)
		if err != nil {
			return nil, err
		}

		result.Orders[key]

		items, err := pis.i.GetItemsWithMaxOrderLineQuantities(req.WarehouseID, req.DistrictID, key, value)

	}
	return
}

func (pis *popularItemServiceImpl) Print(result *models.PopularItemOutput) {

}
