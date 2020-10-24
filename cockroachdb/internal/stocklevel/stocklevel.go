package stocklevel

import (
	"database/sql"
	"fmt"
	"log"
)

// ProcessTransaction processes the Stock Level Transaction
func ProcessTransaction(db *sql.DB, warehouseID, districtID, threshold, lastNOrders int) {
	var totalItems int
	var lastOrderID int

	row := db.QueryRow("SELECT d_next_o_id FROM district WHERE d_w_id=$1 AND d_id=$2", warehouseID, districtID)

	err := row.Scan(&lastOrderID)
	if err != nil {
		log.Fatalf("%v", err)
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
	err = row.Scan(&totalItems)
	if err != nil {
		log.Fatalf("%v", err)
	}

	printOutputState(totalItems, lastOrderID-lastNOrders, lastOrderID)

}

func printOutputState(totalItems, startOrderID, endOrderID int) {
	fmt.Println(fmt.Sprintf("Total Number of Items below threshold: %d , for Order IDs between %d - %d", totalItems, startOrderID, endOrderID))
}
