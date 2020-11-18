package dao

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dbdatamodel"
)

// OrderLineDao provides the interface to interact with the OrderLine table
type OrderLineDao interface {
	Insert(tx *sql.Tx, warehouseID, districtID, orderID int, orderLineItems map[int]*dbdatamodel.OrderLineItem) error
	GetMaxQuantityOrderLinesPerOrder(warehouseID, districtID, startOrderID, lastOrderID int) (map[int]int, error)
	GetOrderLinesForOrder(warehouseID, districtID, orderID int) (map[int]*dbdatamodel.OrderLineItem, error)
	GetDistinctItemIDsPerOrder(warehouseID, districtID int) (map[int][]int, error)
	GetFinalState() (int, error)
}

type orderLineDaoImpl struct {
	db *sql.DB
}

// CreateOrderLineDao creates the dao for order line table
func CreateOrderLineDao(db *sql.DB) OrderLineDao {
	return &orderLineDaoImpl{db: db}
}

// Insert inserts new Order Lines into the table
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

// GetMaxQuantityOrderLinesPerOrder returns the max quantity ordered for items in per order
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

// GetOrderLinesForOrder returns the order lines for the order
func (ol *orderLineDaoImpl) GetOrderLinesForOrder(warehouseID, districtID, orderID int) (result map[int]*dbdatamodel.OrderLineItem, err error) {
	sqlStatement := fmt.Sprintf("SELECT OL_I_ID, OL_SUPPLY_W_ID, OL_QUANTITY, OL_AMOUNT FROM ORDER_LINE_%d_%d WHERE OL_O_ID=%d", warehouseID, districtID, orderID)

	rows, err := ol.db.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("error occurred in getting the order lines. Err: %v", err)
	}
	defer rows.Close()

	var id, supplier, quantity int
	var amount float64

	for rows.Next() {
		if err := rows.Scan(&id, &supplier, &quantity, &amount); err != nil {
			return nil, fmt.Errorf("error occurred in scanning the order line return results. Err: %v", err)
		}
		result[id] = &dbdatamodel.OrderLineItem{
			SupplierWarehouseID: supplier,
			Quantity:            quantity,
			Amount:              amount,
		}
	}

	return
}

// GetDistinctItemIDsPerOrder returns the distint items per order as a map
func (ol *orderLineDaoImpl) GetDistinctItemIDsPerOrder(warehouseID, districtID int) (result map[int][]int, err error) {
	sqlStatement := fmt.Sprintf("SELECT OL_O_ID, OL_I_ID FROM ORDER_LINE_%d_%d GROUP BY OL_O_ID, OL_I_ID ORDER BY OL_O_ID, OL_I_ID", warehouseID, districtID)

	rows, err := ol.db.Query(sqlStatement)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var orderID, itemID int
	for rows.Next() {
		if err := rows.Scan(&orderID, &itemID); err != nil {
			return nil, err
		}

		if result[orderID] == nil {
			result[orderID] = []int{}
		}

		result[orderID] = append(result[orderID], itemID)
	}

	return
}

// GetFinalState returns the final state of the orderline table
func (ol *orderLineDaoImpl) GetFinalState() (totalQuantity int, err error) {
	var tempTotalQuantity int

	baseSQLStatement := "SELECT SUM(OL_QUANTITY) FROM Order_Line_WID_DID"

	for w := 1; w <= 10; w++ {
		for d := 1; d <= 10; d++ {
			finalSQLStatement := strings.ReplaceAll(baseSQLStatement, "WID", strconv.Itoa(w))
			finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "DID", strconv.Itoa(d))

			row := ol.db.QueryRow(finalSQLStatement)
			if err := row.Scan(&tempTotalQuantity); err != nil {
				return 0, err
			}

			totalQuantity += tempTotalQuantity
		}
	}

	return
}
