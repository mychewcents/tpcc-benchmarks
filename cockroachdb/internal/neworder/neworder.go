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
	remote             int
}

// ProcessTransaction process the new transaction
func ProcessTransaction(db *sql.DB, warehouseID, districtID, customerID, numItems int, itemIDs, supplierWarehouseIDs, itemQtys []int) {

	orderTable := fmt.Sprintf("ORDER_%d_%d", warehouseID, districtID)
	// orderLineTable := fmt.Sprintf("ORDER_LINE_%d_%d", warehouseID, districtID)

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
		var remote int

		if supplierWarehouseIDs[index] != warehouseID {
			allLocalFlag = 0
			remote = 1
		}

		itemsToAttributesMap[value] = itemObject{
			supplier: supplierWarehouseIDs[index],
			quantity: itemQtys[index],
			remote:   remote,
		}

		if !itemVisited[value] {
			itemVisited[value] = true
			totalUniqueItems++
		}
	}

	var totalAmount float64

	err = crdb.ExecuteTx(context.Background(), db, nil, func(tx *sql.Tx) error {

		sqlStatement = fmt.Sprintf("INSERT INTO %s (O_ID, O_D_ID, O_W_ID, O_C_ID, O_OL_CNT, O_ALL_LOCAL) VALUES ($1, $2, $3, $5, $6)", orderTable)

		if _, err := tx.Exec(sqlStatement, newOrderID, districtID, warehouseID, customerID, totalUniqueItems, allLocalFlag); err != nil {
			return err
		}

		for key, value := range itemsToAttributesMap {
			// Update Stock

			var itemCurrQty int
			row = tx.QueryRow("SELECT S_QUANTITY, S_YTD FROM STOCK WHERE S_W_ID = $1 AND S_I_ID = $2", value.supplier, key)
			err = row.Scan(&itemCurrQty)

			adjustedQty := itemCurrQty - value.quantity

			if adjustedQty < 10 {
				adjustedQty += 100
			}

			// S_W_ID, S_I_ID, S_QUANTITY, S_YTD, S_ORDER_CNT, S_REMOTE_CNT
			row = tx.QueryRow("UPDATE STOCK SET S_QUANTITY = $1, S_YTD = S_YTD + $2, S_ORDER_CNT = S_ORDER_CNT + 1, S_REMOTE_CNT = $3 WHERE S_W_ID = $4 AND S_I_ID = $5 RETURNING S_I_NAME, S_I_PRICE",
				adjustedQty,
				value.quantity,
				value.remote,
				value.supplier,
				key,
			)

			var itemName string
			var itemPrice float64

			err = row.Scan(&itemName, &itemPrice)

			totalAmount += itemPrice * float64(value.quantity)

			// Add a new Order Line Item String
		}

		return nil
	})
	if err != nil {
		log.Fatalf("%v", err)
	}

	printOutputState(totalAmount)
}

func printOutputState(totalAmount float64) {
	fmt.Println(totalAmount)
}
