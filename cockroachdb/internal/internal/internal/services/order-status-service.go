package services

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dao"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
)

// OrderStatusService interface to processing the transactions
type OrderStatusService interface {
	ProcessTransaction(req *models.OrderStatus) (result *models.OrderStatusOutput, err error)
}

type orderStatusServiceImpl struct {
	db *sql.DB
	o  dao.OrderDao
	c  dao.CustomerDao
	ol dao.OrderLineDao
}

// CreateOrderStatusService creates the service for the Order Status Transaction
func CreateOrderStatusService(db *sql.DB) OrderStatusService {
	return &orderStatusServiceImpl{
		db: db,
		o:  dao.CreateOrderDao(db),
		c:  dao.CreateCustomerDao(db),
		ol: dao.CreateOrderLineDao(db),
	}
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

	orderDetails, err := oss.o.GetLastOrderDetails(req.WarehouseID, req.DistrictID, req.CustomerID)
	if err != nil {
		return nil, err
	}

	customer, err := oss.c.GetDetails(req.WarehouseID, req.DistrictID, req.CustomerID)
	if err != nil {
		return nil, err
	}

	orderLines, err := oss.ol.GetOrderLinesForOrder(req.WarehouseID, req.DistrictID, orderDetails.ID)
	if err != nil {
		return nil, err
	}

	result = &models.OrderStatusOutput{
		Order:      orderDetails,
		Customer:   customer,
		OrderLines: orderLines,
	}

	return nil, nil
}

func (oss *orderStatusServiceImpl) Print(result *models.OrderStatusOutput) {
	var orderStatus strings.Builder

	orderStatus.WriteString(fmt.Sprintf("Customer name: %s %s %s \nBalance: %f", result.Customer.FirstName, result.Customer.MiddleName, result.Customer.LastName, result.Customer.Balance))
	orderStatus.WriteString(fmt.Sprintf("\nOrder: %d \nEntry date: %s \nCarrier: %d", result.Order.ID, result.Order.Timestamp, result.Order.CarrierID))

	for key, value := range result.OrderLines {
		orderStatus.WriteString(fmt.Sprintf("\nItem ID: %d, Supplier: %d, Quantity: %d, Amount: %0.2f", key, value.SupplierWarehouseID, value.Quantity, value.Amount))
	}

	fmt.Println(orderStatus.String())
}
