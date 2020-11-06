package popularitem

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type details struct {
	orderTimestamp         string
	cFirst, cLast, cMiddle string
	maxQuantity            int
	itemNames              []string
}

type itemPercentageName struct {
	percentage float64
	name       string
}

// ProcessTransaction returns the list of the most popular items and their percentage
func ProcessTransaction(db *sql.DB, scanner *bufio.Scanner, transactionArgs []string) bool {
	warehouseID, _ := strconv.Atoi(transactionArgs[0])
	districtID, _ := strconv.Atoi(transactionArgs[1])
	lastNOrders, _ := strconv.Atoi(transactionArgs[2])

	log.Printf("Starting the Popular Item Transaction for: w=%d d=%d n=%d", warehouseID, districtID, lastNOrders)

	if err := execute(db, warehouseID, districtID, lastNOrders); err != nil {
		log.Fatalf("error occurred while executing the popular item transaction. Err: %v", err)
		return false
	}

	// log.Printf("Completed the Popular Item Transaction for: w=%d d=%d n=%d", warehouseID, districtID, lastNOrders)
	return true
}

func execute(db *sql.DB, warehouseID, districtID, lastNOrders int) error {
	// log.Printf("Executing the transaction with the input data...")
	var lastOrderID, startOrderID int

	orderTable := fmt.Sprintf("ORDERS_%d_%d", warehouseID, districtID)
	orderLineTable := fmt.Sprintf("ORDER_LINE_%d_%d", warehouseID, districtID)

	row := db.QueryRow("SELECT d_next_o_id FROM district WHERE d_w_id=$1 AND d_id=$2", warehouseID, districtID)

	if err := row.Scan(&lastOrderID); err != nil {
		return fmt.Errorf("error occured in getting the last order id. Err: %v", err)
	}

	startOrderID = lastOrderID - lastNOrders

	sqlStatement := fmt.Sprintf(`
		SELECT OL_O_ID, MAX(OL_QUANTITY) 
		FROM %s 
		WHERE OL_O_ID < %d 
		AND OL_O_ID >= %d 
		GROUP BY OL_O_ID`,
		orderLineTable, lastOrderID, startOrderID,
	)

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return fmt.Errorf("error occured in getting the maximum order line quantity. Err: %v", err)
	}
	defer rows.Close()

	ordersMap := make(map[int]details)
	itemOccurranceMap := make(map[int]int)
	itemOccurrancePercentageMap := make(map[int]itemPercentageName)

	for rows.Next() {
		// For each order with a maximum quantity of an Order Line item:
		var customerID, orderID, maxQuantity int
		var cFirst, cMiddle, cLast, orderTimestamp string

		if err = rows.Scan(&orderID, &maxQuantity); err != nil {
			return fmt.Errorf("error occured in reading the max order line quantity. Err: %v", err)
		}

		sqlStatement = fmt.Sprintf("SELECT O_C_ID, O_ENTRY_D FROM %s WHERE O_ID = %d", orderTable, orderID)

		row = db.QueryRow(sqlStatement)
		if err = row.Scan(&customerID, &orderTimestamp); err != nil {
			return fmt.Errorf("error occured in getting the details from orders table. Err: %v", err)
		}

		// Fetch the Customer Information
		sqlStatement = fmt.Sprintf("SELECT C_FIRST, C_MIDDLE, C_LAST FROM CUSTOMER WHERE C_W_ID=%d AND C_D_ID = %d AND C_ID = %d", warehouseID, districtID, customerID)

		row = db.QueryRow(sqlStatement)

		if err = row.Scan(&cFirst, &cMiddle, &cLast); err != nil {
			return fmt.Errorf("error occured in getting the customer details. Err: %v", err)
		}

		// Fetch the Item Information
		sqlStatement = fmt.Sprintf("SELECT I_ID, I_NAME FROM ITEM WHERE I_ID IN (SELECT OL_I_ID FROM %s WHERE OL_O_ID = %d AND OL_QUANTITY = %d)", orderLineTable, orderID, maxQuantity)

		items, err := db.Query(sqlStatement)
		if err != nil {
			return fmt.Errorf("error occurred in getting the item details. Err: %v", err)
		}
		defer items.Close()

		var itemIDs []int
		var itemNames []string

		for items.Next() {
			var id int
			var name string

			if err = items.Scan(&id, &name); err != nil {
				return fmt.Errorf("error occurred in scanning the item details. Err: %v", err)
			}

			itemIDs, itemNames = append(itemIDs, id), append(itemNames, name)
			itemOccurranceMap[id]++
		}

		// Calculate the Percentage of orders each items occurred in
		for key, value := range itemIDs {
			percentage := float64((itemOccurranceMap[value] / lastNOrders)) * 100

			itemOccurrancePercentageMap[value] = itemPercentageName{percentage: percentage, name: itemNames[key]}
		}

		ordersMap[orderID] = details{
			maxQuantity:    maxQuantity,
			cFirst:         cFirst,
			cMiddle:        cMiddle,
			cLast:          cLast,
			itemNames:      itemNames,
			orderTimestamp: orderTimestamp,
		}
	}

	// printOutputState(warehouseID, districtID, startOrderID, lastOrderID, lastNOrder, ordersMap, itemOccurrancePercentageMap)
	// log.Printf("Completed executing the transaction with the input data...")
	return nil
}

func printOutputState(warehouseID, districtID, startOrderID, lastOrderID, lastNOrder int, ordersMap map[int]details, itemOccurrancePercentageMap map[int]itemPercentageName) {
	var ordersString strings.Builder

	for key, value := range ordersMap {
		ordersString.WriteString(fmt.Sprintf("\nOrder ID: %d , Timestamp: %s , Max Quantity: %d", key, value.orderTimestamp, value.maxQuantity))
		ordersString.WriteString(fmt.Sprintf("\nCustomer: %s %s %s", value.cFirst, value.cMiddle, value.cLast))
		ordersString.WriteString(fmt.Sprintf("\nItems Ordered: %s", strings.Join(value.itemNames, ", ")))
	}

	var finalPercentageString strings.Builder

	for key, value := range itemOccurrancePercentageMap {
		finalPercentageString.WriteString(fmt.Sprintf("\nItem ID: %d , Name: %s , Percentage: %0.6f", key, value.name, value.percentage))
	}

	fmt.Println(
		fmt.Sprintf(`
			WarehouseID: %d , DistrictID: %d 
			OrderID -> Start: %d , End: %d , Total: %d
			Orders
			%s
			Items Percentage Ratio
			%s
		`,
			warehouseID, districtID, startOrderID, lastOrderID, lastNOrder, ordersString.String(), finalPercentageString.String(),
		),
	)
}
