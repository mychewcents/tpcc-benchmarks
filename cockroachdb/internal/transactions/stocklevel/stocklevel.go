package stocklevel

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

// ProcessTransaction processes the Stock Level Transaction
func ProcessTransaction(db *sql.DB, scanner *bufio.Scanner, transactionArgs []string) bool {
	warehouseID, _ := strconv.Atoi(transactionArgs[0])
	districtID, _ := strconv.Atoi(transactionArgs[1])
	threshold, _ := strconv.Atoi(transactionArgs[2])
	lastNOrders, _ := strconv.Atoi(transactionArgs[3])

	log.Printf("Starting the Stock Level Transaction for: w=%d d=%d t=%d n=%d", warehouseID, districtID, threshold, lastNOrders)
	if err := execute(db, warehouseID, districtID, threshold, lastNOrders); err != nil {
		log.Fatalf("error occurred while executing stock level transaction. Err: %v", err)
		return false
	}

	// log.Printf("Completed the Stock Level Transaction for: w=%d d=%d t=%d n=%d", warehouseID, districtID, threshold, lastNOrders)
	return true
}

func execute(db *sql.DB, warehouseID, districtID, threshold, lastNOrders int) error {
	// log.Printf("Executing the transaction with the input data...")
	var totalItems, lastOrderID int

	row := db.QueryRow("SELECT d_next_o_id FROM district WHERE d_w_id=$1 AND d_id=$2", warehouseID, districtID)

	if err := row.Scan(&lastOrderID); err != nil {
		return fmt.Errorf("error occurred in getting the next order id for the district. Err: %v", err)
	}

	startOrderID := lastOrderID - lastNOrders

	sqlStatement := fmt.Sprintf(`
		SELECT COUNT(*) FROM Stock 
		WHERE S_W_ID=%d 
		AND S_QUANTITY < %d 
		AND S_I_ID IN (
			SELECT OL_I_ID FROM ORDER_LINE_%d_%d 
			WHERE OL_O_ID < %d AND OL_O_ID >= %d
		)`,
		warehouseID, threshold, warehouseID, districtID, lastOrderID, startOrderID,
	)

	row = db.QueryRow(sqlStatement)
	if err := row.Scan(&totalItems); err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error occured in scanning the total items. Err: %v", err)
	}

	// printOutputState(totalItems, lastOrderID-lastNOrders, lastOrderID)
	// log.Printf("Completed executing the transaction with the input data...")
	return nil
}

func printOutputState(totalItems, startOrderID, endOrderID int) {
	fmt.Println(fmt.Sprintf("Total Number of Items below threshold: %d , for Order IDs between %d - %d", totalItems, startOrderID, endOrderID))
}
