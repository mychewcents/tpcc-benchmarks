package dao

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dbdatamodel"
)

// OrderLineDao provides the interface to interact with the OrderLine table
type OrderLineDao interface {
	Insert(tx *sql.Tx, warehouseID, districtID, orderID int, orderLineItems map[int]*dbdatamodel.OrderLineItem) error
	GetMaxQuantityOrderLinesPerOrder(warehouseID, districtID, startOrderID, lastOrderID int) (result map[int]int, err error)
}

type orderLineDaoImpl struct {
	db *sql.DB
}

// CreateOrderLineDao creates the dao for order line table
func CreateOrderLineDao(db *sql.DB) OrderLineDao {
	return &orderLineDaoImpl{db: db}
}

func (ol *orderLineDaoImpl) Insert(tx *sql.Tx, warehouseID, districtID, orderID int, orderLineItems map[int]*dbdatamodel.OrderLineItem) error {
	var orderLineEntries strings.Builder

	idx := 0
	for key, value := range orderLineItems {
		orderLineEntries.WriteString(
			fmt.Sprintf("(%d, %d, %d, %d, %d, %d, %d, %0.2f, '%s'),",
				orderID,
				districtID,
				warehouseID,
				idx+1,
				key,
				value.SupplierWarehouseID,
				value.Quantity,
				value.Amount,
				value.Data,
			))
	}

	orderLineEntriesString := orderLineEntries.String()
	orderLineEntriesString = orderLineEntriesString[:len(orderLineEntriesString)-1]

	orderLineUpdateStatement := fmt.Sprintf("INSERT INTO ORDER_LINE_%d_%d (OL_O_ID, OL_D_ID, OL_W_ID, OL_NUMBER, OL_I_ID, OL_SUPPLY_W_ID, OL_QUANTITY, OL_AMOUNT, OL_DIST_INFO) VALUES %s",
		warehouseID, districtID, orderLineEntriesString)

	if _, err := tx.Exec(orderLineUpdateStatement); err != nil {
		return err
	}

	return nil
}

func (ol *orderLineDaoImpl) GetMaxQuantityOrderLinesPerOrder(warehouseID, districtID, startOrderID, lastOrderID int) (result map[int]int, err error) {
	sqlStatement := fmt.Sprintf(`
		SELECT OL_O_ID, MAX(OL_QUANTITY) 
		FROM ORDER_LINE_%d_%d 
		WHERE OL_O_ID < %d 
		AND OL_O_ID >= %d 
		GROUP BY OL_O_ID`,
		warehouseID, districtID, lastOrderID, startOrderID,
	)

	rows, err := ol.db.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("error occured in getting the maximum order line quantity. Err: %v", err)
	}
	defer rows.Close()

	var orderID, maxQuantity int
	for rows.Next() {
		if err = rows.Scan(&orderID, &maxQuantity); err != nil {
			return nil, fmt.Errorf("error occured in reading the max order line quantity. Err: %v", err)
		}

		result[orderID] = maxQuantity
	}

	return
}
