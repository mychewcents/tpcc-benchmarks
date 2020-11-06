package orderstatus

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/cockroachdb/cockroach-go/crdb"
)

type orderItem struct {
	itemID, supplyWHNumber, quantity int
	itemAmount                       float64
}

// ProcessTransaction processes the Order Status Transaction
func ProcessTransaction(db *sql.DB, scanner *bufio.Scanner, transactionArgs []string) bool {
	warehouseID, _ := strconv.Atoi(transactionArgs[0])
	districtID, _ := strconv.Atoi(transactionArgs[1])
	customerID, _ := strconv.Atoi(transactionArgs[2])

	log.Printf("Starting the Order Status Transaction for: w=%d d=%d c=%d", warehouseID, districtID, customerID)

	if err := execute(db, warehouseID, districtID, customerID); err != nil {
		log.Println("error occured while executing the order status transaction. Err: %v", err)
		return false
	}

	// log.Printf("Completed the Order Status Transaction for: w=%d d=%d c=%d", warehouseID, districtID, customerID)
	return true
}

func execute(db *sql.DB, warehouseID int, districtID int, customerID int) error {
	// log.Printf("Executing the transaction with the input data...")

	lastOrderQuery := fmt.Sprintf("SELECT O_ID, O_DELIVERY_D, O_ENTRY_D, O_CARRIER_ID FROM ORDERS_%d_%d WHERE O_C_ID=%d ORDER BY O_ID DESC LIMIT 1",
		warehouseID, districtID, customerID)
	customerExistsQuery := fmt.Sprintf("SELECT C_FIRST, C_MIDDLE, C_LAST, C_BALANCE  FROM CUSTOMER WHERE C_W_ID=%d AND C_D_ID=%d AND C_ID=%d", warehouseID, districtID, customerID)
	orderLinesQuery := "SELECT OL_I_ID, OL_SUPPLY_W_ID, OL_QUANTITY, OL_AMOUNT FROM ORDER_LINE_%d_%d WHERE OL_O_ID=%d"

	var first, middle, last string
	var balance float64
	if err := db.QueryRow(customerExistsQuery).Scan(&first, &middle, &last, &balance); err != nil {
		return fmt.Errorf("error occured in getting the customers. Err: %v", err)
	}

	var orderLines []orderItem
	var lastOrderID int
	var carrierID sql.NullInt32
	var deliveryDate, entryDate sql.NullString
	err := crdb.ExecuteTx(context.Background(), db, nil, func(tx *sql.Tx) error {
		if err := tx.QueryRow(lastOrderQuery).Scan(&lastOrderID, &deliveryDate, &entryDate, &carrierID); err != nil {
			if err == sql.ErrNoRows {
				log.Printf("no rows for the customer found")
				return nil
			}
			return fmt.Errorf("error occurred in getting the customers. Err: %v", err)
		}
		rows, err := db.Query(fmt.Sprintf(orderLinesQuery, warehouseID, districtID, lastOrderID))
		if err != nil {
			return fmt.Errorf("error occurred in getting the order lines. Err: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var orderLine orderItem
			if err := rows.Scan(&orderLine.itemID, &orderLine.supplyWHNumber, &orderLine.quantity, &orderLine.itemAmount); err != nil {
				return fmt.Errorf("error occurred in scanning the order line return results. Err: %v", err)
			}
			orderLines = append(orderLines, orderLine)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error occured while reading the tables. Err: %v", err)
	}

	// output := "Customer name: %s %s %s \nBalance: %f\nOrder: %d\nEntry date: %s\nCarrier: %s\n"

	// var toPrintDeliveryDate string
	// if !deliveryDate.Valid {
	// 	toPrintDeliveryDate = "undelivered"
	// } else {
	// 	toPrintDeliveryDate = deliveryDate.String
	// }

	// for _, orderLine := range orderLines {

	// 	output += fmt.Sprintf("Item number: %d\nSupply Warehouse: %d\nQuantity: %d\nTotal price: %f\nDelivery date: %s\n",
	// 		orderLine.itemID, orderLine.supplyWHNumber, orderLine.quantity, orderLine.itemAmount, toPrintDeliveryDate)
	// }

	// if carrierID.Valid && carrierID.Int32 != 0 {
	// 	fmt.Println(fmt.Sprintf(output, first, middle, last, balance, lastOrderID, entryDate.String, strconv.Itoa(int(carrierID.Int32))))
	// } else {
	// 	fmt.Println(fmt.Sprintf(output, first, middle, last, balance, lastOrderID, entryDate.String, "null"))
	// }

	// log.Printf("Completed executing the transaction with the input data...")
	return nil
}
