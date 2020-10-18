package neworder

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

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
	var orderLineEntries []string

	err = crdb.ExecuteTx(context.Background(), db, nil, func(tx *sql.Tx) error {

		orderLineItemCounter := 1

		for key, value := range itemsToAttributesMap {
			// Update Stock

			var itemName string
			var itemPrice float64
			var itemCurrQty int
			var itemDistrictData string

			sqlStatement = fmt.Sprintf("SELECT S_I_NAME, S_I_PRICE, S_QUANTITY, S_YTD, S_DIST_%d FROM STOCK WHERE S_W_ID = %d AND S_I_ID = %d", districtID, value.supplier, key)
			row = tx.QueryRow(sqlStatement)
			err = row.Scan(&itemName, &itemPrice, &itemCurrQty, &itemDistrictData)

			adjustedQty := itemCurrQty - value.quantity

			if adjustedQty < 10 {
				adjustedQty += 100
			}

			_, err = tx.Exec("UPDATE STOCK SET S_QUANTITY = $1, S_YTD = S_YTD + $2, S_ORDER_CNT = S_ORDER_CNT + 1, S_REMOTE_CNT = $3 WHERE S_W_ID = $4 AND S_I_ID = $5",
				adjustedQty,
				value.quantity,
				value.remote,
				value.supplier,
				key,
			)

			orderLineAmount := itemPrice * float64(value.quantity)
			totalAmount += orderLineAmount

			// Add a new Order Line Item String
			// OL_O_ID, OL_D_ID, OL_W_ID, OL_NUMBER, OL_I_ID, OL_SUPPLIER_W_ID, OL_QUANTITY, OL_AMOUNT, OL_DIST_INFO
			orderLineEntries = append(orderLineEntries,
				fmt.Sprintf("(%d, %d, %d, %d, %d, %d, %d, %0.2f, %s)",
					newOrderID,
					districtID,
					warehouseID,
					orderLineItemCounter,
					key,
					value.supplier,
					value.quantity,
					orderLineAmount,
					itemDistrictData,
				))

			orderLineItemCounter++
		}

		sqlStatement = fmt.Sprintf("INSERT INTO %s (O_ID, O_D_ID, O_W_ID, O_C_ID, O_OL_CNT, O_ALL_LOCAL, O_TOTAL_AMOUNT) VALUES ($1, $2, $3, $5, $6, $7)", orderTable)

		if _, err := tx.Exec(sqlStatement, newOrderID, districtID, warehouseID, customerID, totalUniqueItems, allLocalFlag, totalAmount); err != nil {
			return err
		}

		sqlStatement = fmt.Sprintf("INSERT INTO %s (OL_O_ID, OL_D_ID, OL_W_ID, OL_NUMBER, OL_I_ID, OL_SUPPLY_W_ID, OL_QUANTITY, OL_AMOUNT, OL_DIST_INFO) VALUES %s",
			orderLineTable, strings.Join(orderLineEntries, ", "))

		if _, err := tx.Exec(sqlStatement); err != nil {
			return err
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
