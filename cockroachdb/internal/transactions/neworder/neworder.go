package neworder

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/cockroachdb/cockroach-go/crdb"
)

type itemObject struct {
	id           int
	quantity     int
	supplier     int
	remote       int
	startStock   int
	finalStock   int
	currYTD      float64
	currOrderCnt int
	name         string
	price        float64
	data         string
}

// ProcessTransaction process the new order transaction
func ProcessTransaction(db *sql.DB, scanner *bufio.Scanner, args []string) bool {

	customerID, _ := strconv.Atoi(args[0])
	warehouseID, _ := strconv.Atoi(args[1])
	districtID, _ := strconv.Atoi(args[2])
	numOfItems, _ := strconv.Atoi(args[3])

	log.Printf("Starting the New Order Transaction for row: c=%d w=%d d=%d n=%d", customerID, warehouseID, districtID, numOfItems)

	log.Printf("Pre-processing the input data...")
	orderLineObjects := make([]*itemObject, numOfItems)
	orderItems := make([]int, numOfItems)

	prevSeenOrderLineItems := make(map[int]int) // Maps the items IDs to the
	isLocal := 1

	var id, supplier, quantity, remote, totalUniqueItems int

	for i := 0; i < numOfItems; i++ {
		if scanner.Scan() {
			olArgs := strings.Split(scanner.Text(), ",")
			id, _ = strconv.Atoi(olArgs[0])
			supplier, _ = strconv.Atoi(olArgs[1])
			quantity, _ = strconv.Atoi(olArgs[2])

			if supplier != warehouseID {
				remote = 1
				if isLocal == 1 {
					isLocal = 0
				}
			} else {
				remote = 0
			}

			if _, ok := prevSeenOrderLineItems[id]; ok {
				orderLineObjects[prevSeenOrderLineItems[id]].quantity += quantity
			} else {
				orderLineObjects[totalUniqueItems] = &itemObject{
					id:       id,
					supplier: supplier,
					quantity: quantity,
					remote:   remote,
				}
				prevSeenOrderLineItems[id] = totalUniqueItems
				orderItems[totalUniqueItems] = id
				totalUniqueItems++
			}
		}
	}

	orderLineObjects = orderLineObjects[0:totalUniqueItems]
	orderItems = orderItems[0:totalUniqueItems]
	sort.Ints(orderItems)

	log.Printf("Completed pre-processing the input data...")

	if err := execute(db, warehouseID, districtID, customerID, totalUniqueItems, isLocal, totalUniqueItems, orderLineObjects, orderItems); err != nil {
		log.Fatalf("error occured while executing the new order transaction. Err: %v", err)
		return false
	}

	log.Printf("Completed the New Order Transaction for row: c=%d w=%d d=%d n=%d", customerID, warehouseID, districtID, numOfItems)
	return true
}

func execute(db *sql.DB, warehouseID, districtID, customerID, numItems, isLocal, totalUniqueItems int, orderLineObjects []*itemObject, sortedOrderItems []int) error {
	// log.Printf("Executing the transaction with the input data...")

	orderTable := fmt.Sprintf("ORDERS_%d_%d", warehouseID, districtID)
	orderLineTable := fmt.Sprintf("ORDER_LINE_%d_%d", warehouseID, districtID)
	orderItemCustomerPairTable := fmt.Sprintf("ORDER_ITEMS_CUSTOMERS_%d_%d", warehouseID, districtID)

	sqlStatement := fmt.Sprintf("UPDATE District SET D_NEXT_O_ID = D_NEXT_O_ID + 1 WHERE D_W_ID = %d AND D_ID = %d RETURNING D_NEXT_O_ID, D_TAX, D_W_TAX", warehouseID, districtID)

	var newOrderID int
	var districtTax, warehouseTax float64
	row := db.QueryRow(sqlStatement)
	if err := row.Scan(&newOrderID, &districtTax, &warehouseTax); err != nil {
		return fmt.Errorf("error occured in updating the district table for the next order id. Err: %v", err)
	}

	var cLastName, cCredit string
	var cDiscount float64
	sqlStatement = fmt.Sprintf("SELECT C_LAST, C_CREDIT, C_DISCOUNT FROM CUSTOMER WHERE C_W_ID=%d AND C_D_ID = %d AND C_ID = %d", warehouseID, districtID, customerID)

	row = db.QueryRow(sqlStatement)

	if err := row.Scan(&cLastName, &cCredit, &cDiscount); err != nil {
		return fmt.Errorf("error occured in getting the customer details. Err: %v", err)
	}

	var totalAmount float64
	var orderTimestamp string

	var orderItemCustomerPair strings.Builder

	for i := 0; i < len(sortedOrderItems)-1; i++ {
		for j := i + 1; j < len(sortedOrderItems); j++ {
			orderItemCustomerPair.WriteString(fmt.Sprintf("(%d, %d, %d, %d, %d),", warehouseID, districtID, customerID, sortedOrderItems[i], sortedOrderItems[j]))
		}
	}

	// var ch chan bool
	err := crdb.ExecuteTx(context.Background(), db, nil, func(tx *sql.Tx) error {
		var orderLineEntries []string

		// log.Printf("Starting the insert of the Item Pair for the customer: c=%d w=%d d=%d ", customerID, warehouseID, districtID)
		insertItemPairsParallel(tx, orderItemCustomerPairTable, orderItemCustomerPair.String())

		for key, value := range orderLineObjects {
			sqlStatement = fmt.Sprintf("SELECT S_I_NAME, S_I_PRICE, S_QUANTITY, S_YTD, S_ORDER_CNT, S_DIST_%02d FROM STOCK WHERE S_W_ID = %d AND S_I_ID = %d", districtID, value.supplier, value.id)
			row = tx.QueryRow(sqlStatement)
			if err := row.Scan(&value.name, &value.price, &value.startStock, &value.currYTD, &value.currOrderCnt, &value.data); err != nil {
				return fmt.Errorf("error in getting the stock details for the item: w=%d id=%d \nErr: %v", value.supplier, value.id, err)
			}

			adjustedQty := value.startStock - value.quantity
			if adjustedQty < 10 {
				adjustedQty += 100
			}
			value.finalStock = adjustedQty

			orderLineAmount := value.price * float64(value.quantity)
			totalAmount += orderLineAmount

			// Add a new Order Line Item String
			// OL_O_ID, OL_D_ID, OL_W_ID, OL_NUMBER, OL_I_ID, OL_SUPPLIER_W_ID, OL_QUANTITY, OL_AMOUNT, OL_DIST_INFO
			orderLineEntries = append(orderLineEntries,
				fmt.Sprintf("(%d, %d, %d, %d, %d, %d, %d, %0.2f, '%s')",
					newOrderID,
					districtID,
					warehouseID,
					key,
					value.id,
					value.supplier,
					value.quantity,
					orderLineAmount,
					value.data,
				))
		}

		bulkUpdatesOrderLineItems := make([]string, totalUniqueItems)
		bulkQuantityUpdates := make([]string, totalUniqueItems)
		bulkYTDUpdates := make([]string, totalUniqueItems)
		bulkOrderCountUpdates := make([]string, totalUniqueItems)
		bulkRemoteCountUpdates := make([]string, totalUniqueItems)

		idx := 0
		for _, value := range orderLineObjects {
			bulkUpdatesOrderLineItems[idx] = fmt.Sprintf("(%d, %d)", value.supplier, value.id)
			bulkQuantityUpdates[idx] = fmt.Sprintf("WHEN (%d, %d) THEN %d", value.supplier, value.id, value.finalStock)
			bulkYTDUpdates[idx] = fmt.Sprintf("WHEN (%d, %d) THEN %d", value.supplier, value.id, int(value.currYTD)+value.quantity)
			bulkOrderCountUpdates[idx] = fmt.Sprintf("WHEN (%d, %d) THEN %d", value.supplier, value.id, value.currOrderCnt+1)
			bulkRemoteCountUpdates[idx] = fmt.Sprintf("WHEN (%d, %d) THEN %d", value.supplier, value.id, value.remote)
			idx++
		}

		sqlStatement = fmt.Sprintf(`
			UPDATE STOCK 
				SET S_QUANTITY = CASE (S_W_ID, S_I_ID) %s END, 
				S_YTD = CASE (S_W_ID, S_I_ID) %s END, 
				S_ORDER_CNT = CASE (S_W_ID, S_I_ID) %s END, 
				S_REMOTE_CNT = CASE (S_W_ID, S_I_ID) %s END 
			WHERE (S_W_ID, S_I_ID) IN (%s)`,
			strings.Join(bulkQuantityUpdates, " "),
			strings.Join(bulkYTDUpdates, " "),
			strings.Join(bulkOrderCountUpdates, " "),
			strings.Join(bulkRemoteCountUpdates, " "),
			strings.Join(bulkUpdatesOrderLineItems, ", "),
		)

		if _, err := tx.Exec(sqlStatement); err != nil {
			return fmt.Errorf("error in updating the stock details \nErr: %v", err)
		}

		sqlStatement = fmt.Sprintf("INSERT INTO %s (O_ID, O_D_ID, O_W_ID, O_C_ID, O_OL_CNT, O_ALL_LOCAL, O_TOTAL_AMOUNT) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING O_ENTRY_D", orderTable)

		totalAmount = totalAmount * (1.0 + districtTax + warehouseTax) * (1.0 - cDiscount)

		row = tx.QueryRow(sqlStatement, newOrderID, districtID, warehouseID, customerID, totalUniqueItems, isLocal, totalAmount)
		if err := row.Scan(&orderTimestamp); err != nil {
			return fmt.Errorf("error in inserting new order row: w=%d d=%d o=%d \n Err: %v", warehouseID, districtID, newOrderID, err)
		}

		sqlStatement = fmt.Sprintf("INSERT INTO %s (OL_O_ID, OL_D_ID, OL_W_ID, OL_NUMBER, OL_I_ID, OL_SUPPLY_W_ID, OL_QUANTITY, OL_AMOUNT, OL_DIST_INFO) VALUES %s",
			orderLineTable, strings.Join(orderLineEntries, ", "))

		if _, err := tx.Exec(sqlStatement); err != nil {
			return fmt.Errorf("error in inserting new order line rows: w=%d d=%d o=%d \n Err: %v", warehouseID, districtID, newOrderID, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error occured in updating the order/order lines/item pairs table. Err: %v", err)
	}

	// printOutputState(warehouseID, districtID, customerID, cLastName, cCredit, cDiscount,
	// 	newOrderID, orderTimestamp, totalUniqueItems, totalAmount, orderLineObjects)
	// log.Printf("Completed executing the transaction with the input data...")
	return nil
}

func insertItemPairsParallel(tx *sql.Tx, orderItemCustomerPairTable string, orderItemCustomerPair string) {
	// log.Printf("Inserting the item pairs")
	sqlStatement := fmt.Sprintf("UPSERT INTO %s (IC_W_ID, IC_D_ID, IC_C_ID, IC_I_1_ID, IC_I_2_ID) VALUES %s", orderItemCustomerPairTable, orderItemCustomerPair)
	sqlStatement = sqlStatement[0 : len(sqlStatement)-1]

	if _, err := tx.Exec(sqlStatement); err != nil {
		// ch <- false
		log.Fatalf("error occured in the item pairs for customers. Err: %v", err)
	}

	// log.Printf("Completed inserting the item pairs")
	// ch <- true
}

func printOutputState(warehouseID, districtID, customerID int, cLastName, cCredit string, cDiscount float64,
	orderID int, orderTimestamp string, totalUniqueItems int, totalAmount float64, orderLineObjects []*itemObject) {
	var newOrderString strings.Builder

	newOrderString.WriteString(fmt.Sprintf("Customer Identifier => W_ID = %d, D_ID = %d, C_ID = %d \n", warehouseID, districtID, customerID))
	newOrderString.WriteString(fmt.Sprintf("Customer Info => Last Name: %s , Credit: %s , Discount: %0.6f \n", cLastName, cCredit, cDiscount))
	newOrderString.WriteString(fmt.Sprintf("Order Details: O_ID = %d , O_ENTRY_D = %s \n", orderID, orderTimestamp))
	newOrderString.WriteString(fmt.Sprintf("Total Unique Items: %d \n", totalUniqueItems))
	newOrderString.WriteString(fmt.Sprintf("Total Amount: %.2f \n", totalAmount))

	newOrderString.WriteString(fmt.Sprintf(" # \t ID \t Name (Supplier, Qty, Amount, Stock) \n"))
	for key, value := range orderLineObjects {
		newOrderString.WriteString(fmt.Sprintf(" %02d \t %d \t %s (%d, %d, %.2f, %d) \n",
			key+1,
			value.id,
			value.name,
			value.supplier,
			value.quantity,
			value.price*float64(value.quantity),
			value.finalStock,
		))
	}

	fmt.Println(newOrderString.String())
}
