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

	return execute(db, warehouseID, districtID, customerID)
}

func execute(db *sql.DB, warehouseID, districtID, customerID int) bool {

	relatedCustomerIdentifiers := make(map[int]map[int]map[int]bool)
	orderItemCustomerPairTable := "ORDER_ITEMS_CUSTOMERS_WID_DID"

	var orderLineItemPairString strings.Builder

	sqlStatement := fmt.Sprintf("SELECT IC_I_1_ID, IC_I_2_ID FROM ORDER_ITEMS_CUSTOMERS_%d_%d WHERE IC_C_ID = %d", warehouseID, districtID, customerID)

	rows, err := db.Query(sqlStatement)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no rows found for customer: %d %d %d", warehouseID, districtID, customerID)
			return true
		}
		log.Fatalf("error in fetching the order line item pairs for the asked customer: %d %d %d \n Err: %v", warehouseID, districtID, customerID, err)
		return false
	}

	var itemID1, itemID2 int
	for rows.Next() {
		err := rows.Scan(&itemID1, &itemID2)
		if err != nil {
			log.Fatalf("error in fetching the reading the order line item pair: %d %d %d \n Err: %v", warehouseID, districtID, customerID, err)
		}
		orderLineItemPairString.WriteString(fmt.Sprintf("(IC_I_1_ID = %d AND IC_I_2_ID = %d) OR ", itemID1, itemID2))
	}

	finalOrderLineItemPairWhereClause := orderLineItemPairString.String()

	if len(finalOrderLineItemPairWhereClause) == 0 {
		log.Fatalf("could not create the final WHERE clause script for related customer: %d %d %d", warehouseID, districtID, customerID)
		return false
	}

	finalOrderLineItemPairWhereClause = finalOrderLineItemPairWhereClause[:len(finalOrderLineItemPairWhereClause)-4]

	baseSQLStatement := fmt.Sprintf(`
		SELECT IC_C_ID FROM %s p WHERE %s
	`, orderItemCustomerPairTable, finalOrderLineItemPairWhereClause)

	var cCustomerID int

	for w := 1; w <= 10; w++ {
		if w != warehouseID {
			for d := 1; d <= 10; d++ {
				finalSQLStatement := strings.ReplaceAll(baseSQLStatement, "WID", strconv.Itoa(w))
				finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "DID", strconv.Itoa(d))

				rows, err := db.Query(finalSQLStatement)
				if err == sql.ErrNoRows {
					continue
				}
				if err != nil {
					fmt.Println(err)
					return false
				}

				for rows.Next() {
					err := rows.Scan(&cCustomerID)
					if err != nil {
						fmt.Println(err)
					}
					if !relatedCustomerIdentifiers[w][d][cCustomerID] {

						if relatedCustomerIdentifiers[w] == nil {
							relatedCustomerIdentifiers[w] = make(map[int]map[int]bool)
						}
						if relatedCustomerIdentifiers[w][d] == nil {
							relatedCustomerIdentifiers[w][d] = make(map[int]bool)
						}
						relatedCustomerIdentifiers[w][d][cCustomerID] = true
					}
				}
			}

		}
	}

	// printOutputState(warehouseID, districtID, customerID, relatedCustomerIdentifiers)
	return true
}

func printOutputState(warehouseID, districtID, customerID int, relatedCustomerIdentifiers map[int]map[int]map[int]bool) {
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
