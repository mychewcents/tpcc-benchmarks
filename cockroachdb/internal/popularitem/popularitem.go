package populatitem

import (
	"database/sql"
	"fmt"
	"log"
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
func ProcessTransaction(db *sql.DB, warehouseID, districtID, lastNOrder int) {
	var lastOrderID, startOrderID int

	orderTable := fmt.Sprintf("ORDER_%d_%d", warehouseID, districtID)
	orderLineTable := fmt.Sprintf("ORDER_LINE_%d_%d", warehouseID, districtID)

	row := db.QueryRow(`SELECT d_next_o_id FROM district WHERE d_w_id=$1 AND d_id=$2`, warehouseID, districtID)

	err := row.Scan(&lastOrderID)
	if err != nil {
		log.Fatalf("%v", err)
	}

	startOrderID = lastOrderID - lastNOrder

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
		log.Fatalf("%v", err)
	}
	defer rows.Close()

	ordersMap := make(map[int]details)
	itemOccurranceMap := make(map[int]int)
	itemOccurrancePercentageMap := make(map[int]itemPercentageName)

	for rows.Next() {
		// For each order with a maximum quantity of an Order Line item:
		var orderID, maxQuantity int
		err = rows.Scan(&orderID, &maxQuantity)
		if err != nil {
			log.Fatalf("%v", err)
		}

		sqlStatement = fmt.Sprintf(`SELECT O_ENTRY_D FROM %s WHERE O_ID = %d`, orderTable, orderID)

		var orderTimestamp string
		row = db.QueryRow(sqlStatement)
		err = row.Scan(&orderTimestamp)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Fetch the Customer Information
		sqlStatement = fmt.Sprintf(`SELECT C_FIRST, C_MIDDLE, C_LAST FROM CUSTOMERS WHERE C_W_ID=%d AND C_D_ID = %d AND C_ID = %d`, warehouseID, districtID, orderID)

		var cFirst, cMiddle, cLast string
		row = db.QueryRow(sqlStatement)
		err = row.Scan(&cFirst, &cMiddle, &cLast)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Fetch the Item Information
		sqlStatement = fmt.Sprintf(`SELECT I_ID, I_NAME FROM ITEM WHERE I_ID IN (SELECT OL_I_ID FROM %s WHERE OL_O_ID = %d AND OL_QUANTITY = %d`, orderLineTable, orderID, maxQuantity)

		items, err := db.Query(sqlStatement)
		if err != nil {
			log.Fatalf("%v", err)
		}
		defer items.Close()

		var itemIDs []int
		var itemNames []string

		for items.Next() {
			var id int
			var name string
			err = items.Scan(&id, &name)
			if err != nil {
				log.Fatalf("%v", err)
			}

			itemIDs, itemNames = append(itemIDs, id), append(itemNames, name)
			itemOccurranceMap[id]++
		}

		// Calculate the Percentage of orders each items occurred in
		for key := range itemIDs {
			percentage := float64((itemOccurranceMap[key] / lastNOrder)) * 100

			itemOccurrancePercentageMap[key] = itemPercentageName{percentage: percentage, name: itemNames[key]}
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

	outputState(warehouseID, districtID, startOrderID, lastOrderID, lastNOrder, ordersMap, itemOccurrancePercentageMap)
}

func outputState(warehouseID, districtID, startOrderID, lastOrderID, lastNOrder int, ordersMap map[int]details, itemOccurrancePercentageMap map[int]itemPercentageName) {
	var ordersString strings.Builder

	for key, value := range ordersMap {
		ordersString.WriteString(fmt.Sprintf("\nOrder ID: %d , Timestamp: %s , Max Quantity: %d", key, value.orderTimestamp, value.maxQuantity))
		ordersString.WriteString(fmt.Sprintf("\nCustomer: %s %s %s", value.cFirst, value.cMiddle, value.cLast))
		ordersString.WriteString(fmt.Sprintf("\nItems Ordered: %s", strings.Join(value.itemNames, ", ")))
	}

	var finalPercentageString strings.Builder

	for key, value := range itemOccurrancePercentageMap {
		finalPercentageString.WriteString(fmt.Sprintf("\nItem ID: %d , Name: %s , Percentage: %0.02f", key, value.name, value.percentage))
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
