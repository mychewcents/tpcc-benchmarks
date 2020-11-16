package dao

import (
	"database/sql"
	"fmt"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dbdatamodel"
)

// OrderDao interface to the Order Partitioned Table
type OrderDao interface {
	Insert(tx *sql.Tx, warehouseID, districtID, customerID, orderID, uniqueItems, isOrderLocal int, totalAmount float64) (string, error)
	GetDetails(warehouseID, districtID, orderID int) (*dbdatamodel.Order, error)
	GetLastOrderDetails(warehouseID, districtID, customerID int) (*dbdatamodel.Order, error)
	GetUndeliveredOrderIDsPerDistrict(warehouseID int) (map[int]int, error)
	DeliverOrder(tx *sql.Tx, warehouseID, districtID, orderID, carrierID int) (int, float64, error)
	GetOrderIDForCustomers(warehouseID, districtID int) (map[int]int, error)
}

type orderDaoImpl struct {
	db *sql.DB
}

// CreateOrderDao creates new object for Orders table
func CreateOrderDao(db *sql.DB) OrderDao {
	return &orderDaoImpl{db: db}
}

// Insert adds new order row to the database
func (od *orderDaoImpl) Insert(tx *sql.Tx, warehouseID, districtID, customerID, orderID, uniqueItems, isOrderLocal int, totalAmount float64) (orderTimestamp string, err error) {

	orderUpdateStatement := fmt.Sprintf(`
		INSERT INTO ORDERS_%d_%d (O_ID, O_D_ID, O_W_ID, O_C_ID, O_OL_CNT, O_ALL_LOCAL, O_TOTAL_AMOUNT) 
		VALUES (%d, %d, %d, %d, %d, %d, %0.2f) 
		RETURNING O_ENTRY_D`,
		warehouseID,
		districtID,
		orderID,
		warehouseID,
		districtID,
		customerID,
		uniqueItems,
		isOrderLocal,
		totalAmount,
	)

	row := tx.QueryRow(orderUpdateStatement)
	if err := row.Scan(&orderTimestamp); err != nil {
		return orderTimestamp, err
	}

	return orderTimestamp, nil
}

// GetDetails returns the order details
func (od *orderDaoImpl) GetDetails(warehouseID, districtID, orderID int) (result *dbdatamodel.Order, err error) {
	sqlStatement := fmt.Sprintf("SELECT O_C_ID, O_ENTRY_D FROM ORDERS_%d_%d WHERE O_ID = %d", warehouseID, districtID, orderID)

	var customerID int
	var orderTimestamp string

	row := od.db.QueryRow(sqlStatement)
	if err = row.Scan(&customerID, &orderTimestamp); err != nil {
		return nil, fmt.Errorf("error occured in getting the details from orders table. Err: %v", err)
	}

	result = &dbdatamodel.Order{
		ID:          orderID,
		WarehouseID: warehouseID,
		DistrictID:  districtID,
		CustomerID:  customerID,
		Timestamp:   orderTimestamp,
	}

	return
}

// GetLastOrderDetails returns the details of the last order placed by the customer
func (od *orderDaoImpl) GetLastOrderDetails(warehouseID, districtID, customerID int) (result *dbdatamodel.Order, err error) {
	sqlStatement := fmt.Sprintf("SELECT O_ID, O_DELIVERY_D, O_ENTRY_D, O_CARRIER_ID FROM ORDERS_%d_%d WHERE O_C_ID=%d ORDER BY O_ID DESC LIMIT 1",
		warehouseID, districtID, customerID)

	var orderID, carrierID int
	var orderTimestamp string
	var deliveryTimestamp sql.NullString

	row := od.db.QueryRow(sqlStatement)
	if err = row.Scan(&orderID, &deliveryTimestamp, &orderTimestamp, &carrierID); err != nil {
		return nil, fmt.Errorf("error occured in getting the details from orders table. Err: %v", err)
	}

	result = &dbdatamodel.Order{
		ID:          orderID,
		WarehouseID: warehouseID,
		DistrictID:  districtID,
		CustomerID:  customerID,
		Timestamp:   orderTimestamp,
		CarrierID:   carrierID,
	}

	if deliveryTimestamp.Valid {
		result.DeliveryTimestamp = deliveryTimestamp.String
	}

	return
}

// GetUndeliveredOrderIDsPerDistrict returns the last undelivered order ids for each of the districts of the warehouse
func (od *orderDaoImpl) GetUndeliveredOrderIDsPerDistrict(warehouseID int) (result map[int]int, err error) {
	baseSQLStatement := "SELECT O_ID FROM ORDERS_%d_%d WHERE O_CARRIER_ID=0 ORDER BY O_ID LIMIT 1"

	var orderID int

	for dID := 1; dID <= 10; dID++ {
		finalSQLStatement := fmt.Sprintf(baseSQLStatement, warehouseID, dID)

		err := od.db.QueryRow(finalSQLStatement).Scan(&orderID)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			} else {
				return nil, fmt.Errorf("error occured while fetching the orders. Err: %v", err)
			}
		}
		result[dID] = orderID
	}

	return
}

// DeliverOrder updates the order after delivery
func (od *orderDaoImpl) DeliverOrder(tx *sql.Tx, warehouseID, districtID, orderID, carrierID int) (customerID int, totalAmount float64, err error) {
	panic("implement me")
}

// GetOrderLineCountForCustomer returns the order line count for each customer
func (od *orderDaoImpl) GetOrderIDForCustomers(warehouseID, districtID int) (countPerCustomer map[int]int, err error) {

	sqlStatement := fmt.Sprintf("SELECT O_ID, O_C_ID FROM ORDERS_%d_%d", warehouseID, districtID)

	rows, err := od.db.Query(sqlStatement)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var orderID, customerID int
	for rows.Next() {
		if err := rows.Scan(&orderID, &customerID); err != nil {
			return nil, fmt.Errorf("error in getting the order id for w = %d d = %d", warehouseID, districtID)
		}
		countPerCustomer[customerID] = orderID
	}

	return
}
