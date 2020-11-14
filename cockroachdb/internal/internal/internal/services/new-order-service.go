package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/cockroachdb/cockroach-go/crdb"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dao"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
)

// NewOrderService provides the implementation of the New Order transaciton
type NewOrderService interface {
	ProcessNewOrderTransaction(req *models.NewOrder) (*models.NewOrderOutput, error)
}

// NewOrderServiceImpl stores the new order input and output models
type NewOrderServiceImpl struct {
	d   dao.DistrictDao
	c   dao.CustomerDao
	s   dao.StockDao
	o   dao.OrderDao
	ol  dao.OrderLineDao
	cip dao.CustomerItemsPairDao
	db  *sql.DB
}

// GetNewOrderService returns the object for a new order transaction
func GetNewOrderService(db *sql.DB) NewOrderService {
	return &NewOrderServiceImpl{
		db:  db,
		d:   dao.CreateDistrictDao(db),
		c:   dao.CreateCustomerDao(db),
		s:   dao.CreateStockDao(db),
		o:   dao.CreateOrderDao(db),
		ol:  dao.CreateOrderLineDao(db),
		cip: dao.CreateCustomerItemsPairDao(db),
	}
}

// ProcessNewOrderTransaction process the new order transaction
func (nos *NewOrderServiceImpl) ProcessNewOrderTransaction(req *models.NewOrder) (*models.NewOrderOutput, error) {
	log.Printf("Starting the New Order Transaction for row: c=%d w=%d d=%d n=%d", req.CustomerID, req.WarehouseID, req.DistrictID, req.UniqueItems)

	result, err := nos.execute(req)
	if err != nil {
		return nil, fmt.Errorf("error occured while executing the new order transaction. Err: %v", err)
	}

	log.Printf("Completed the New Order Transaction for row: c=%d w=%d d=%d n=%d", req.CustomerID, req.WarehouseID, req.DistrictID, req.UniqueItems)
	return result, nil
}

func (nos *NewOrderServiceImpl) execute(req *models.NewOrder) (*models.NewOrderOutput, error) {
	// log.Printf("Executing the transaction with the input data...")

	newOrderID, districtTax, warehouseTax, err := nos.d.GetNewOrderIDAndTaxRates(req.WarehouseID, req.DistrictID)
	if err != nil {
		return nil, err
	}
	result := &models.NewOrderOutput{
		DistrictTax:  districtTax,
		WarehouseTax: warehouseTax,
		OrderID:      newOrderID,
	}

	customer, err := nos.c.GetCustomerDetails(req.WarehouseID, req.DistrictID, req.CustomerID)
	if err != nil {
		return nil, err
	}
	result.Customer = customer

	if err := nos.cip.Insert(req.WarehouseID, req.DistrictID, req.CustomerID, req.UniqueItems, req.NewOrderLineItems); err != nil {
		return nil, err
	}

	if err := crdb.ExecuteTx(context.Background(), nos.db, nil, func(tx *sql.Tx) error {

		totalAmount, err := nos.s.GetStockDetails(tx, req.DistrictID, req.NewOrderLineItems)
		if err != nil {
			return err
		}

		result.TotalOrderAmount = totalAmount

		if err := nos.s.UpdateStockDetails(tx, req.NewOrderLineItems); err != nil {
			return fmt.Errorf("error in updating stock table: w=%d d=%d o=%d \n Err: %v", req.WarehouseID, req.DistrictID, newOrderID, err)
		}

		result.OrderTimestamp, err = nos.o.Insert(tx, req.WarehouseID, req.DistrictID, req.CustomerID, newOrderID, req.UniqueItems, req.IsOrderLocal, req.TotalAmount)
		if err != nil {
			return fmt.Errorf("error in inserting new order row: w=%d d=%d o=%d \n Err: %v", req.WarehouseID, req.DistrictID, newOrderID, err)
		}

		if err := nos.ol.Insert(tx, req.WarehouseID, req.DistrictID, newOrderID, req.NewOrderLineItems); err != nil {
			return fmt.Errorf("error in inserting new order line rows: w=%d d=%d o=%d \n Err: %v", req.WarehouseID, req.DistrictID, newOrderID, err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error occured in updating the order/order lines/item pairs table. Err: %v", err)
	}

	result.TotalOrderAmount = req.TotalAmount * (1.0 + result.DistrictTax + result.WarehouseTax) * (1.0 - result.Customer.Discount)
	result.UniqueItems = req.UniqueItems
	result.OrderLineItems = req.NewOrderLineItems

	// nos.Print(result)

	return result, nil
}

// Print prints the formatted output of the NewOrder Transaction
func (nos *NewOrderServiceImpl) Print(o *models.NewOrderOutput) {
	var newOrderString strings.Builder

	newOrderString.WriteString(fmt.Sprintf("Customer Identifier => W_ID = %d, D_ID = %d, C_ID = %d \n", o.Customer.WarehouseID, o.Customer.DistrictID, o.Customer.CustomerID))
	newOrderString.WriteString(fmt.Sprintf("Customer Info => Last Name: %s , Credit: %s , Discount: %0.6f \n", o.Customer.LastName, o.Customer.Credit, o.Customer.Discount))
	newOrderString.WriteString(fmt.Sprintf("Order Details: O_ID = %d , O_ENTRY_D = %s \n", o.OrderID, o.OrderTimestamp))
	newOrderString.WriteString(fmt.Sprintf("Total Unique Items: %d \n", o.UniqueItems))
	newOrderString.WriteString(fmt.Sprintf("Total Amount: %.2f \n", o.TotalOrderAmount))

	newOrderString.WriteString(fmt.Sprintf(" # \t ID \t Name (Supplier, Qty, Amount, Stock) \n"))
	idx := 1
	for key, value := range o.OrderLineItems {
		newOrderString.WriteString(fmt.Sprintf(" %02d \t %d \t %s (%d, %d, %.2f, %d) \n",
			idx,
			key,
			value.Name,
			value.SupplierWarehouseID,
			value.Quantity,
			value.Price*float64(value.Quantity),
			value.FinalStock,
		))
		idx++
	}

	fmt.Println(newOrderString.String())
}
