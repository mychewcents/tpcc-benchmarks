package neworder

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/cockroachdb/cockroach-go/crdb"
)

type newOrder struct {
	// INPUTS
	CustomerID, DistrictID, WarehouseID, NumItems int
	ItemIDs, SupplierWarehouseIDs, ItemQuantities []int64

	// OUTPUTS
	lastName, creditStatus, orderTimestamp               string
	custDiscount, totalAmount, warehouseTax, districtTax float64
	orderID                                              int
	itemNames                                            []string
	itemAmount                                           []float64
	itemStock                                            []int
}

type itemObject struct {
	quantity, supplier int
}

// ProcessTransaction process the new transaction
func ProcessTransaction(db *sql.DB, warehouseID, districtID, customerID, numItems int, itemIDs, supplierWarehouseIDs, itemQtys []int) {

	orderTable := fmt.Sprintf("ORDER_%d_%d", warehouseID, districtID)
	orderLineTable := fmt.Sprintf("ORDER_LINE_%d_%d", warehouseID, districtID)

	sqlStatement := fmt.Sprintf("UPDATE District SET D_NEXT_O_ID = D_NEXT_O_ID + 1 WHERE D_W_ID = %d AND D_ID = %d RETURNING D_NEXT_O_ID, D_TAX, D_W_TAX", warehouseID, districtID)

	var newOrderID int
	var districtTax, warehouseTax float64
	row := db.QueryRow(sqlStatement)
	err := row.Scan(&newOrderID, &districtTax, &warehouseTax)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	var totalUniqueItems int

	allLocalFlag := 1
	itemVisited := make(map[int]bool)
	itemsToAttributesMap := make(map[int]itemObject)

	for index, value := range itemIDs {
		itemsToAttributesMap[value] = itemObject{
			supplier: supplierWarehouseIDs[index],
			quantity: itemQtys[index],
		}
		if !itemVisited[value] {
			itemVisited[value] = true
			totalUniqueItems++
		}
		if allLocalFlag == 1 && supplierWarehouseIDs[index] != warehouseID {
			allLocalFlag = 0
		}
	}

	err = crdb.ExecuteTx(context.Background(), db, nil, func(tx *sql.Tx) error {

		sqlStatement = fmt.Sprintf("INSERT INTO %s (O_ID, O_D_ID, O_W_ID, O_C_ID, O_OL_CNT, O_ALL_LOCAL) VALUES ($1, $2, $3, $5, $6)", orderTable)

		if _, err := tx.Exec(sqlStatement, newOrderID, districtID, warehouseID, customerID, totalUniqueItems, allLocalFlag); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatalf("%v", err)
	}

	printOutputState()
}

func printOutputState() {
	fmt.Println()
}
