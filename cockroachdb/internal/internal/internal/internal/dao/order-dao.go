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
