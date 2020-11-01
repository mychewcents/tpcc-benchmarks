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

	printOutputState(warehouseID, districtID, customerID)
	return true
}

func printOutputState(warehouseID, districtID, customerID int) {
	var relatedCustomerString strings.Builder

	fmt.Println(relatedCustomerString.String())
}
