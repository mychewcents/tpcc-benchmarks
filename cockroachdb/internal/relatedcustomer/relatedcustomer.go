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
	orderItemCustomerPairTable := "ORDER_ITEMS_CUSTOMERS_WID"

	sqlStatement := fmt.Sprintf(`
		SELECT p.IC_W_ID, p.IC_D_ID, p.IC_C_ID 
		FROM %s p 
		INNER JOIN 
		(SELECT * FROM ORDER_ITEMS_CUSTOMERS_%d WHERE IC_W_ID = %d AND IC_D_ID = %d AND IC_C_ID = %d) c 
		ON p.IC_I_1_ID = c.IC_I_1_ID AND p.IC_I_2_ID = c.IC_I_2_ID
	`, orderItemCustomerPairTable, warehouseID, warehouseID, districtID, customerID)

	var cWarehouseID, cDistrictID, cCustomerID int

	for i := 1; i < 11; i++ {
		if i != warehouseID {
			rows, err := db.Query(strings.ReplaceAll(sqlStatement, "WID", strconv.Itoa(i)))
			if err == sql.ErrNoRows {
				continue
			}
			if err != nil {
				fmt.Println(err)
				return false
			}

			for rows.Next() {
				err := rows.Scan(&cWarehouseID, &cDistrictID, &cCustomerID)
				if err != nil {
					fmt.Println(err)
				}
				if !relatedCustomerIdentifiers[cWarehouseID][cDistrictID][cCustomerID] {

					if relatedCustomerIdentifiers[cWarehouseID] == nil {
						relatedCustomerIdentifiers[cWarehouseID] = make(map[int]map[int]bool)
					}
					if relatedCustomerIdentifiers[cWarehouseID][cDistrictID] == nil {
						relatedCustomerIdentifiers[cWarehouseID][cDistrictID] = make(map[int]bool)
					}
					relatedCustomerIdentifiers[cWarehouseID][cDistrictID][cCustomerID] = true
				}

			}
		}
	}

	printOutputState(warehouseID, districtID, customerID, relatedCustomerIdentifiers)
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
