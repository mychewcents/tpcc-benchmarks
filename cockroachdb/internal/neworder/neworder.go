package neworder

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
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
	id, quantity, supplier, remote int
}

// ProcessTransaction process the new order transaction
func ProcessTransaction(db *sql.DB, scanner *bufio.Scanner, args []string) {
	customerID, _ := strconv.Atoi(args[0])
	warehouseID, _ := strconv.Atoi(args[1])
	districtID, _ := strconv.Atoi(args[2])
	numOfItems, _ := strconv.Atoi(args[3])

	orderLineObjects := make([]*itemObject, numOfItems)
	prevSeenOrderLineItems := make(map[int]int) // Maps the items IDs to the
	isLocal := 1

	var id, supplier, quantity, remote, totalUniqueItems int

	for i := 0; i < numOfItems; i++ {
		olArgs := strings.Split(scanner.Text(), ",")
		id, _ = strconv.Atoi(olArgs[0])
		supplier, _ = strconv.Atoi(olArgs[1])
		quantity, _ = strconv.Atoi(olArgs[2])

		if supplier != warehouseID {
			remote = 1
			if isLocal == 1 {
				isLocal = 0
			}
		} else {
			remote = 0
		}

		if _, ok := prevSeenOrderLineItems[id]; ok {
			orderLineObjects[prevSeenOrderLineItems[id]].quantity += quantity
		} else {
			orderLineObjects[i] = &itemObject{
				id:       id,
				supplier: supplier,
				quantity: quantity,
				remote:   remote,
			}
			prevSeenOrderLineItems[id] = i
			totalUniqueItems++
		}

	}

	execute(db, warehouseID, districtID, customerID, totalUniqueItems, isLocal, totalUniqueItems, orderLineObjects[0:totalUniqueItems])
}

func execute(db *sql.DB, warehouseID, districtID, customerID, numItems, isLocal, totalUniqueItems int, orderLineObjects []*itemObject) {

	orderTable := fmt.Sprintf("ORDERS_%d_%d", warehouseID, districtID)
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

	var totalAmount float64
	var orderLineEntries []string

	err = crdb.ExecuteTx(context.Background(), db, nil, func(tx *sql.Tx) error {

		var itemName string
		var itemPrice float64
		var itemCurrQty int
		var itemDistrictData string

		for key, value := range orderLineObjects {

			sqlStatement = fmt.Sprintf("SELECT S_I_NAME, S_I_PRICE, S_QUANTITY, S_DIST_%02d FROM STOCK WHERE S_W_ID = %d AND S_I_ID = %d", districtID, value.supplier, value.id)
			row = tx.QueryRow(sqlStatement)
			err = row.Scan(&itemName, &itemPrice, &itemCurrQty, &itemDistrictData)
			if err != nil {
				return err
			}

			// fmt.Println(itemName, itemPrice, itemCurrQty, itemDistrictData)
			adjustedQty := itemCurrQty - value.quantity

			if adjustedQty < 10 {
				adjustedQty += 100
			}

			sqlStatement = fmt.Sprintf("UPDATE STOCK SET S_QUANTITY = %d, S_YTD = S_YTD + %d, S_ORDER_CNT = S_ORDER_CNT + 1, S_REMOTE_CNT = %d WHERE S_W_ID = %d AND S_I_ID = %d",
				adjustedQty,
				value.quantity,
				value.remote,
				value.supplier,
				value.id)
			_, err = tx.Exec(sqlStatement)

			// fmt.Println(sqlStatement)

			_, err = tx.Exec(sqlStatement)
			orderLineAmount := itemPrice * float64(value.quantity)
			totalAmount += orderLineAmount

			// Add a new Order Line Item String
			// OL_O_ID, OL_D_ID, OL_W_ID, OL_NUMBER, OL_I_ID, OL_SUPPLIER_W_ID, OL_QUANTITY, OL_AMOUNT, OL_DIST_INFO
			orderLineEntries = append(orderLineEntries,
				fmt.Sprintf("(%d, %d, %d, %d, %d, %d, %d, %0.2f, '%s')",
					newOrderID,
					districtID,
					warehouseID,
					key,
					value.id,
					value.supplier,
					value.quantity,
					orderLineAmount,
					itemDistrictData,
				))

		}

		sqlStatement = fmt.Sprintf("INSERT INTO %s (O_ID, O_D_ID, O_W_ID, O_C_ID, O_OL_CNT, O_ALL_LOCAL, O_TOTAL_AMOUNT) VALUES ($1, $2, $3, $4, $5, $6, $7)", orderTable)

		if _, err := tx.Exec(sqlStatement, newOrderID, districtID, warehouseID, customerID, totalUniqueItems, isLocal, totalAmount); err != nil {
			return err
		}

		sqlStatement = fmt.Sprintf("INSERT INTO %s (OL_O_ID, OL_D_ID, OL_W_ID, OL_NUMBER, OL_I_ID, OL_SUPPLY_W_ID, OL_QUANTITY, OL_AMOUNT, OL_DIST_INFO) VALUES %s",
			orderLineTable, strings.Join(orderLineEntries, ", "))

		// fmt.Println(sqlStatement)
		if _, err := tx.Exec(sqlStatement); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatalf("%v", err)
	}

	printOutputState(warehouseID, districtID, customerID, totalAmount)
}

func printOutputState(warehouseID, districtID, customerID int, totalAmount float64) {
	fmt.Println(totalAmount)
}
