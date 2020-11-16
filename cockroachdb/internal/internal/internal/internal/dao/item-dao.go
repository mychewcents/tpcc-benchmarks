package dao

import (
	"database/sql"
	"fmt"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dbdatamodel"
)

// ItemDao interface to the Item table queries
type ItemDao interface {
	GetItemsWithMaxOrderLineQuantities(warehouseID, districtID, orderID, maxQuantity int) (map[int]*dbdatamodel.Item, error)
}

type itemDaoImpl struct {
	db *sql.DB
}

// CreateItemDao creates a new object for the Item tables
func CreateItemDao(db *sql.DB) ItemDao {
	return &itemDaoImpl{db: db}
}

func (id *itemDaoImpl) GetItemsWithMaxOrderLineQuantities(warehouseID, districtID, orderID, maxQuantity int) (items map[int]*dbdatamodel.Item, err error) {
	sqlStatement := fmt.Sprintf("SELECT I_ID, I_NAME FROM ITEM WHERE I_ID IN (SELECT OL_I_ID FROM ORDER_LINE_%d_%d WHERE OL_O_ID = %d AND OL_QUANTITY = %d)",
		warehouseID, districtID, orderID, maxQuantity)

	rows, err := id.db.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("error occurred in getting the item details. Err: %v", err)
	}
	defer rows.Close()

	var itemID int
	var name string
	for rows.Next() {
		if err = rows.Scan(&itemID, &name); err != nil {
			return nil, fmt.Errorf("error occurred in scanning the item details. Err: %v", err)
		}

		items[itemID] = &dbdatamodel.Item{ID: itemID, Name: name}
	}

	return
}
