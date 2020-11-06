package relatedcustomer

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// ProcessTransaction processes the Related Customer Transaction
func ProcessTransaction(db *sql.DB, scanner *bufio.Scanner, transactionArgs []string) bool {
	warehouseID, _ := strconv.Atoi(transactionArgs[0])
	districtID, _ := strconv.Atoi(transactionArgs[1])
	customerID, _ := strconv.Atoi(transactionArgs[2])

	log.Printf("Starting the Related Customer Transaction for: w=%d d=%d c=%d", warehouseID, districtID, customerID)

	if err := execute(db, warehouseID, districtID, customerID); err != nil {
		log.Fatalf("error occurred in executing the related customer transaction. Err: %v", err)
		return false
	}

	log.Printf("Completed the Related Customer Transaction for: w=%d d=%d c=%d", warehouseID, districtID, customerID)
	return true
}

func execute(db *sql.DB, warehouseID, districtID, customerID int) error {
	log.Printf("Executing the transaction with the input data...")

	relatedCustomerIdentifiers := make(map[int]map[int][]int)
	orderItemCustomerPairTable := "ORDER_ITEMS_CUSTOMERS_WID_DID"

	var orderLineItemPairString strings.Builder

	sqlStatement := fmt.Sprintf("SELECT IC_I_1_ID, IC_I_2_ID FROM ORDER_ITEMS_CUSTOMERS_%d_%d WHERE IC_C_ID = %d", warehouseID, districtID, customerID)

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return fmt.Errorf("error in fetching the order line item pairs. Err: %v", err)
	}
	defer rows.Close()

	var itemID1, itemID2 int
	for rows.Next() {
		err := rows.Scan(&itemID1, &itemID2)
		if err != nil {
			return fmt.Errorf("error occurred in scanning the order line item pair. Err: %v", err)
		}
		orderLineItemPairString.WriteString(fmt.Sprintf("(%d, %d),", itemID1, itemID2))
	}

	finalOrderLineItemPairWhereClause := orderLineItemPairString.String()

	if len(finalOrderLineItemPairWhereClause) == 0 {
		log.Printf("could not create the final WHERE clause script. No item pairs found")
		return nil
	}

	finalOrderLineItemPairWhereClause = finalOrderLineItemPairWhereClause[:len(finalOrderLineItemPairWhereClause)-1]

	baseSQLStatement := fmt.Sprintf("SELECT IC_C_ID FROM %s p WHERE (IC_I_1_ID, IC_I_2_ID) IN (%s)", orderItemCustomerPairTable, finalOrderLineItemPairWhereClause)

	ch := make(chan []int, 90)

	for w := 1; w <= 10; w++ {
		if w != warehouseID {
			for d := 1; d <= 10; d++ {
				finalSQLStatement := strings.ReplaceAll(baseSQLStatement, "WID", strconv.Itoa(w))
				finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "DID", strconv.Itoa(d))

				go getRelatedCustomersParallel(db, w, d, finalSQLStatement, ch)
			}
		}
	}

	for i := 0; i < 90; i++ {
		relatedCustomersArray := <-ch
		w := relatedCustomersArray[0]
		d := relatedCustomersArray[1]
		log.Println(w, d)
		if relatedCustomerIdentifiers[w] == nil {
			relatedCustomerIdentifiers[w] = make(map[int][]int)
		}

		if len(relatedCustomersArray) > 2 {
			cIDs := relatedCustomersArray[2:]
			relatedCustomerIdentifiers[w][d] = cIDs
		}
	}

	// printOutputState(warehouseID, districtID, customerID, relatedCustomerIdentifiers)
	log.Printf("Completed executing the transaction with the input data...")
	return nil
}

func getRelatedCustomersParallel(db *sql.DB, w, d int, finalSQLStatement string, ch chan []int) {
	var relatedCustomers []int
	relatedCustomers = append(relatedCustomers, w)
	relatedCustomers = append(relatedCustomers, d)

	rows, err := db.Query(finalSQLStatement)
	if err == sql.ErrNoRows {
		ch <- nil
		return
	}
	if err != nil {
		log.Fatalf("error occurred in reading the related customers from table: w=%d d=%d. Err: %v", w, d, err)
		return
	}

	var cID int

	for rows.Next() {
		err := rows.Scan(&cID)
		if err != nil {
			log.Fatalf("error occurred in scanning the related customer id. Err: %v", err)
			return
		}
		relatedCustomers = append(relatedCustomers, cID)
	}

	ch <- relatedCustomers
}

func printOutputState(warehouseID, districtID, customerID int, relatedCustomerIdentifiers map[int]map[int][]int) {
	var relatedCustomerIdentifierString, relatedCustomerString strings.Builder

	for wKey, wValue := range relatedCustomerIdentifiers {
		for dKey, dValue := range wValue {
			for cKey := range dValue {
				relatedCustomerIdentifierString.WriteString(fmt.Sprintf("(%d,%d,%d),", wKey, dKey, cKey))
			}
		}
	}

	finalCustomersString := relatedCustomerIdentifierString.String()
	if len(finalCustomersString) > 0 {
		finalCustomersString = finalCustomersString[0 : len(finalCustomersString)-1]
	}

	relatedCustomerString.WriteString(fmt.Sprintf("Related Customers for the customer: %d, %d, %d \n", warehouseID, districtID, customerID))
	relatedCustomerString.WriteString(fmt.Sprintf("%s", finalCustomersString))

	fmt.Println(relatedCustomerString.String())
}
