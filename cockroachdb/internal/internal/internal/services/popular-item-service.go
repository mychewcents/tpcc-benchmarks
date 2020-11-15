package services

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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

		items, err := pis.i.GetItemsWithMaxOrderLineQuantities(req.WarehouseID, req.DistrictID, key, value)
		if err != nil {
			return nil, err
		}

		result.Orders[key] = &models.PopularItemOrderDetails{Order: orderDetails, Customer: customer, Items: items, MaxOLQuantity: value}
		for key, value := range items {
			if _, ok := result.ItemOccurances[key]; ok {
				result.ItemOccurances[key].Occurances++
			} else {
				result.ItemOccurances[key] = &models.PopularItemOccuranceAndPercentage{Occurances: 1, Name: value.Name}
			}
		}
	}

	for _, value := range result.ItemOccurances {
		value.Percentage = float64((value.Occurances / req.LastNOrders) * 100)
	}

	return
}

func (pis *popularItemServiceImpl) Print(result *models.PopularItemOutput) {
	var popularItemResult strings.Builder
	var ordersResult strings.Builder

	popularItemResult.WriteString(fmt.Sprintf("WarehouseID: %d, DistrictID: %d", result.WarehouseID, result.DistrictID))
	popularItemResult.WriteString(fmt.Sprintf("OrderIDs: Start: %d, End: %d, Total: %d", result.StartOrderID, result.LastOrderID, result.LastOrderID-result.StartOrderID))

	for key, value := range result.Orders {
		ordersResult.WriteString(fmt.Sprintf("\nOrder ID: %d, Timestamp: %s, Max Quantity: %d", key, value.Order.Timestamp, value.MaxOLQuantity))
		ordersResult.WriteString(fmt.Sprintf("\nCustomer: %s %s %s", value.Customer.FirstName, value.Customer.MiddleName, value.Customer.LastName))
		ordersResult.WriteString(fmt.Sprintf("\nItems ordered: "))

		for itemKey, itemValue := range value.Items {
			ordersResult.WriteString(fmt.Sprintf("\nID: %d; Name: %s", itemKey, itemValue.Name))
		}
	}

	ordersResult.WriteString(fmt.Sprintf("\nMost Popular Items: "))
	for key, value := range result.ItemOccurances {
		ordersResult.WriteString(fmt.Sprintf("\nID: %d, Name: %s, Percentage: %0.6f", key, value.Name, value.Percentage))
	}

	popularItemResult.WriteString(fmt.Sprintf("\n%s", ordersResult.String()))

	popularItemResultString := popularItemResult.String()
	fmt.Println(popularItemResultString)
}
