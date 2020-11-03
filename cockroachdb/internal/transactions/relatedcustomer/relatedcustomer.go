package relatedcustomer

import (
	"bufio"
	"database/sql"
	"fmt"
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

	baseSQLStatement := fmt.Sprintf(`
		SELECT p.IC_W_ID, p.IC_D_ID, p.IC_C_ID 
		FROM %s p 
		INNER JOIN 
		(SELECT * FROM ORDER_ITEMS_CUSTOMERS_%d_%d WHERE IC_C_ID = %d) c 
		ON p.IC_I_1_ID = c.IC_I_1_ID AND p.IC_I_2_ID = c.IC_I_2_ID
	`, orderItemCustomerPairTable, warehouseID, districtID, customerID)

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
