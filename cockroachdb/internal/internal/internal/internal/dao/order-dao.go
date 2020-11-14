package dao

import (
	"database/sql"
	"fmt"
)

// OrderDao interface to the Order Partitioned Table
type OrderDao interface {
	Insert(tx *sql.Tx, warehouseID, districtID, customerID, orderID, uniqueItems, isOrderLocal int, totalAmount float64) (string, error)
}

type orderDaoImpl struct {
	db *sql.DB
}

// CreateOrderDao creates new object for Orders table
func CreateOrderDao(db *sql.DB) OrderDao {
	return &orderDaoImpl{db: db}
}

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
