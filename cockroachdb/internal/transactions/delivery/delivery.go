package delivery

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/cockroachdb/cockroach-go/crdb"
)

type districtOrder struct {
	districtID, orderID int
}

// ProcessTransaction processes the Delivery transaction
func ProcessTransaction(db *sql.DB, scanner *bufio.Scanner, transactionArgs []string) bool {
	warehouseID, _ := strconv.Atoi(transactionArgs[0])
	carrierID, _ := strconv.Atoi(transactionArgs[1])

	log.Printf("Starting the Delivery Transaction for: w=%d c=%d", warehouseID, carrierID)

	if err := execute(db, warehouseID, carrierID); err != nil {
		log.Println("error occured while executing the delivery transaction. Err: %v", err)
		return false
	}

	// log.Printf("Completed the Delivery Transaction for: w=%d c=%d", warehouseID, carrierID)
	return true
}

func execute(db *sql.DB, warehouseID int, carrierID int) error {
	// log.Printf("Executing the transaction with the input data...")

	orderQuery := "SELECT O_ID FROM ORDERS_%d_%d WHERE O_CARRIER_ID=0 ORDER BY O_ID LIMIT 1"
	updateOrderQuery := "UPDATE ORDERS_%d_%d SET (O_CARRIER_ID, O_DELIVERY_D) = (%d, now()) WHERE O_W_ID=%d AND O_D_ID=%d AND O_ID=%d RETURNING O_C_ID, O_TOTAL_AMOUNT"
	updateCustomerQuery := "UPDATE CUSTOMER SET (C_BALANCE, C_DELIVERY_CNT) = (C_BALANCE + %f, C_DELIVERY_CNT + 1) WHERE C_W_ID=%d AND C_D_ID=%d AND C_ID=%d"

	var orders []districtOrder

	// var orderIDs [10]int32
	for district := 1; district <= 10; district++ {
		var orderID sql.NullInt32
		if err := db.QueryRow(fmt.Sprintf(orderQuery, warehouseID, district)).Scan(&orderID); err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("error occured while fetching the orders. Err: %v", err)
		}
		if orderID.Valid {
			// orderIDs[district-1] = orderID.Int32
			orders = append(orders, districtOrder{district, int(orderID.Int32)})
		}
	}

	if len(orders) == 0 {
		log.Printf("no pending orders to deliver found.")
		return nil
	}

	err := crdb.ExecuteTx(context.Background(), db, nil, func(tx *sql.Tx) error {
		for _, order := range orders {
			districtID := order.districtID
			orderID := order.orderID
			var totalAmount float64
			var customerID int
			if err := tx.QueryRow(fmt.Sprintf(updateOrderQuery, warehouseID, districtID, carrierID, warehouseID, districtID, orderID)).Scan(&customerID, &totalAmount); err != nil {
				return fmt.Errorf("error occurred while updating the order details. Err: %v", err)
			}
			if _, err := tx.Exec(fmt.Sprintf(updateCustomerQuery, totalAmount, warehouseID, districtID, customerID)); err != nil {
				return fmt.Errorf("error occurred while updating the customer details. Err: %v", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error occurred while updating the order/customer table. Err: %v", err)
	}

	// log.Printf("Completed executing the transaction with the input data...")
	return nil
}
