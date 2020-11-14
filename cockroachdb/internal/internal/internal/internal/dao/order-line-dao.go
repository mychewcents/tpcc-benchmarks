package dao

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
)

// OrderLineDao provides the interface to interact with the OrderLine table
type OrderLineDao interface {
	Insert(tx *sql.Tx, warehouseID, districtID, orderID int, orderLineItems map[int]*models.NewOrderOrderLineItem) error
}

type orderLineDaoImpl struct {
	db *sql.DB
}

// CreateOrderLineDao creates the dao for order line table
func CreateOrderLineDao(db *sql.DB) OrderLineDao {
	return &orderLineDaoImpl{db: db}
}

func (ol *orderLineDaoImpl) Insert(tx *sql.Tx, warehouseID, districtID, orderID int, orderLineItems map[int]*models.NewOrderOrderLineItem) error {
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
