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

// HandleTransaction performs the transaction and returns the execution result in Boolean
func HandleTransaction(db *sql.DB, scanner *bufio.Scanner, args []string, printOutput bool) bool {
	cID, _ := strconv.Atoi(args[0])
	wID, _ := strconv.Atoi(args[1])
	dID, _ := strconv.Atoi(args[2])
	numOfOrderLineItems, _ := strconv.Atoi(args[3])

	newOrderLines, isOrderLocal, totalUniqueItems := readAndPrepareOrderLineItems(scanner, numOfOrderLineItems, wID)

	n := &NewOrder{
		CustomerID:        cID,
		WarehouseID:       wID,
		DistrictID:        dID,
		IsOrderLocal:      isOrderLocal,
		UniqueItems:       totalUniqueItems,
		NewOrderLineItems: newOrderLines,
	}

	return n.ProcessTransaction(db, printOutput)
}

func readAndPrepareOrderLineItems(scanner *bufio.Scanner, numOfItems, warehouseID int) (orderLineItems map[int]*OrderLineItem, isOrderLocal, totalUniqueOrderItems int) {
	orderLineItems = make(map[int]*OrderLineItem)
	isOrderLocal = 1

	var id, supplier, quantity, remote int

	for i := 0; i < numOfItems; i++ {
		if scanner.Scan() {
			args := strings.Split(scanner.Text(), ",")

			id, _ = strconv.Atoi(args[0])
			supplier, _ = strconv.Atoi(args[1])
			quantity, _ = strconv.Atoi(args[2])

			if supplier != warehouseID {
				remote = 1
				if isOrderLocal == 1 {
					isOrderLocal = 0
				}
			} else {
				remote = 0
			}

			if _, ok := orderLineItems[id]; ok {
				orderLineItems[id].Quantity += quantity
			} else {
				orderLineItems[id] = &OrderLineItem{
					SupplierWarehouseID: supplier,
					Quantity:            quantity,
					IsRemote:            remote,
				}
				totalUniqueOrderItems++
			}
		}
	}

	return
}

// ProcessTransaction process the new order transaction
func (n *NewOrder) ProcessTransaction(db *sql.DB, printResult bool) bool {
	log.Printf("Starting the New Order Transaction for row: c=%d w=%d d=%d n=%d", n.CustomerID, n.WarehouseID, n.DistrictID, n.UniqueItems)

	result, err := n.execute(db)
	if err != nil {
		log.Printf("error occured while executing the new order transaction. Err: %v", err)
		return false
	}

	if printResult {
		result.Print()
	}

	log.Printf("Completed the New Order Transaction for row: c=%d w=%d d=%d n=%d", n.CustomerID, n.WarehouseID, n.DistrictID, n.UniqueItems)
	return true
}

func (n *NewOrder) getNewOrderIDAndTaxRates(db *sql.DB) (newOrderID int, wTax, dTax float64, err error) {
	sqlStatement := fmt.Sprintf("UPDATE District SET D_NEXT_O_ID = D_NEXT_O_ID + 1 WHERE D_W_ID = $1 AND D_ID = $2 RETURNING D_NEXT_O_ID, D_TAX, D_W_TAX")

	row := db.QueryRow(sqlStatement, n.WarehouseID, n.DistrictID)
	if err := row.Scan(&newOrderID, &dTax, &wTax); err != nil {
		return 0, 0.0, 0.0, fmt.Errorf("error occured in updating the district table for the next order id. Err: %v", err)
	}

	return
}

func (n *NewOrder) getCustomerInformation(db *sql.DB) (cLastName, cCredit string, cDiscount float64, err error) {
	sqlStatement := fmt.Sprintf("SELECT C_LAST, C_CREDIT, C_DISCOUNT FROM CUSTOMER WHERE C_W_ID = $1 AND C_D_ID = $2 AND C_ID = $3")

	row := db.QueryRow(sqlStatement, n.WarehouseID, n.DistrictID, n.CustomerID)
	if err := row.Scan(&cLastName, &cCredit, &cDiscount); err != nil {
		return "", "", 0.0, fmt.Errorf("error occured in getting the customer details. Err: %v", err)
	}

	return
}

func (n *NewOrder) insertOrderPairItems(db *sql.DB) error {
	var orderItemCustomerPair strings.Builder
	orderItemCustomerPairTable := fmt.Sprintf("ORDER_ITEMS_CUSTOMERS_%d_%d", n.WarehouseID, n.DistrictID)

	sortedOrderItems := make([]int, n.UniqueItems)

	idx := 0
	for key := range n.NewOrderLineItems {
		sortedOrderItems[idx] = key
		idx++
	}

	sort.Ints(sortedOrderItems)

	for i := 0; i < len(sortedOrderItems)-1; i++ {
		for j := i + 1; j < len(sortedOrderItems); j++ {
			orderItemCustomerPair.WriteString(fmt.Sprintf("(%d, %d, %d, %d, %d),", n.WarehouseID, n.DistrictID, n.CustomerID, sortedOrderItems[i], sortedOrderItems[j]))
		}
	}

	sqlStatement := fmt.Sprintf("UPSERT INTO %s (IC_W_ID, IC_D_ID, IC_C_ID, IC_I_1_ID, IC_I_2_ID) VALUES %s", orderItemCustomerPairTable, orderItemCustomerPair)
	sqlStatement = sqlStatement[0 : len(sqlStatement)-1]

	if _, err := db.Exec(sqlStatement); err != nil {
		return fmt.Errorf("error occured in the item pairs for customers. Err: %v", err)
	}

	return nil
}

func (n *NewOrder) getItemDetails(tx *sql.Tx) error {

	var itemsWhereClause strings.Builder

	for key, value := range n.NewOrderLineItems {
		itemsWhereClause.WriteString(fmt.Sprintf("(%d, %d),", value.SupplierWarehouseID, key))
	}

	itemsWhereClauseString := itemsWhereClause.String()
	itemsWhereClauseString = itemsWhereClauseString[:len(itemsWhereClauseString)-1]

	sqlStatement := fmt.Sprintf("SELECT S_I_ID, S_I_NAME, S_I_PRICE, S_QUANTITY, S_YTD, S_ORDER_CNT, S_DIST_%02d FROM STOCK WHERE (S_W_ID, S_I_ID) IN %s", n.DistrictID, itemsWhereClauseString)
	rows, err := tx.Query(sqlStatement)
	if err == sql.ErrNoRows {
		return fmt.Errorf("no rows found for the items ids passed")
	}
	if err != nil {
		return fmt.Errorf("error in getting the stock details for the items. \nquery: %s. \nErr: %v", sqlStatement, err)
	}

	var name, data string
	var price, currYTD float64
	var id, startStock, currOrderCnt int

	for rows.Next() {
		if err := rows.Scan(&id, &name, &price, &startStock, &currYTD, &currOrderCnt, &data); err != nil {
			return fmt.Errorf("error in scanning the results for the items. Err: %v", err)
		}

		if value, ok := n.NewOrderLineItems[id]; ok {
			value.Name = name
			value.Price = price
			value.StartStock = startStock
			value.CurrYTD = currYTD
			value.CurrOrderCnt = currOrderCnt
			value.Data = data

			adjustedQty := startStock - value.Quantity
			if adjustedQty < 10 {
				adjustedQty += 100
			}
			value.FinalStock = adjustedQty

			value.Amount = price * float64(value.Quantity)
			n.TotalAmount += value.Amount
		}
	}

	return nil
}

func (n *NewOrder) prepareStatements(orderID int) (orderUpdateStatement, orderLineUpdateStatement, stockUpdateStatement string) {
	var orderLineEntries, stockOrderItemIdentifiers, stockQuantityUpdates, stockYTDUpdates, stockOrderCountUpdates, stockRemoteCountUpdates strings.Builder

	var itemIdentifier string
	whenClauseFormat := "WHEN %s THEN %d "

	idx := 0
	for key, value := range n.NewOrderLineItems {
		orderLineEntries.WriteString(
			fmt.Sprintf("(%d, %d, %d, %d, %d, %d, %d, %0.2f, '%s'),",
				orderID,
				n.DistrictID,
				n.WarehouseID,
				idx+1,
				key,
				value.SupplierWarehouseID,
				value.Quantity,
				value.Amount,
				value.Data,
			))

		itemIdentifier = fmt.Sprintf("(%d, %d)", value.SupplierWarehouseID, key)

		stockOrderItemIdentifiers.WriteString(fmt.Sprintf("%s,", itemIdentifier))
		stockQuantityUpdates.WriteString(fmt.Sprintf(whenClauseFormat, itemIdentifier, value.FinalStock))
		stockYTDUpdates.WriteString(fmt.Sprintf(whenClauseFormat, itemIdentifier, int(value.CurrYTD)+value.Quantity))
		stockOrderCountUpdates.WriteString(fmt.Sprintf(whenClauseFormat, itemIdentifier, value.CurrOrderCnt+1))
		stockRemoteCountUpdates.WriteString(fmt.Sprintf(whenClauseFormat, itemIdentifier, value.IsRemote))
		idx++
	}

	orderLineEntriesString := orderLineEntries.String()
	orderLineEntriesString = orderLineEntriesString[:len(orderLineEntriesString)-1]

	stockOrderItemIdentifiersString := stockOrderItemIdentifiers.String()
	stockOrderItemIdentifiersString = stockOrderItemIdentifiersString[:len(stockOrderItemIdentifiersString)-1]

	orderUpdateStatement = fmt.Sprintf(`
		INSERT INTO ORDERS_%d_%d (O_ID, O_D_ID, O_W_ID, O_C_ID, O_OL_CNT, O_ALL_LOCAL, O_TOTAL_AMOUNT) 
		VALUES (%d, %d, %d, %d, %d, %d, %0.2f) 
		RETURNING O_ENTRY_D`,
		n.WarehouseID,
		n.DistrictID,
		orderID,
		n.WarehouseID,
		n.DistrictID,
		n.CustomerID,
		n.UniqueItems,
		n.IsOrderLocal,
		n.TotalAmount,
	)

	orderLineUpdateStatement = fmt.Sprintf("INSERT INTO ORDER_LINE_%d_%d (OL_O_ID, OL_D_ID, OL_W_ID, OL_NUMBER, OL_I_ID, OL_SUPPLY_W_ID, OL_QUANTITY, OL_AMOUNT, OL_DIST_INFO) VALUES %s",
		n.WarehouseID, n.DistrictID, orderLineEntriesString)

	stockUpdateStatement = fmt.Sprintf(`
			UPDATE STOCK 
				SET S_QUANTITY = CASE (S_W_ID, S_I_ID) %s END, 
				S_YTD = CASE (S_W_ID, S_I_ID) %s END, 
				S_ORDER_CNT = CASE (S_W_ID, S_I_ID) %s END, 
				S_REMOTE_CNT = CASE (S_W_ID, S_I_ID) %s END 
			WHERE (S_W_ID, S_I_ID) IN (%s)`,
		stockQuantityUpdates.String(),
		stockYTDUpdates.String(),
		stockOrderCountUpdates.String(),
		stockRemoteCountUpdates.String(),
		stockOrderItemIdentifiersString,
	)

	return orderUpdateStatement, orderLineUpdateStatement, stockUpdateStatement
}

func (n *NewOrder) execute(db *sql.DB) (result *Output, err error) {
	// log.Printf("Executing the transaction with the input data...")

	newOrderID, districtTax, warehouseTax, err := n.getNewOrderIDAndTaxRates(db)
	if err != nil {
		return nil, err
	}
	result = &Output{
		Customer: &CustomerInfo{
			WarehouseID: n.WarehouseID,
			DistrictID:  n.DistrictID,
			CustomerID:  n.CustomerID,
		},
		DistrictTax:  districtTax,
		WarehouseTax: warehouseTax,
		OrderID:      newOrderID,
	}

	cLastName, cCredit, cDiscount, err := n.getCustomerInformation(db)
	if err != nil {
		return nil, err
	}
	result.Customer.LastName = cLastName
	result.Customer.Credit = cCredit
	result.Customer.Discount = cDiscount

	if err := n.insertOrderPairItems(db); err != nil {
		return nil, err
	}

	if err := crdb.ExecuteTx(context.Background(), db, nil, func(tx *sql.Tx) error {
		if err := n.getItemDetails(tx); err != nil {
			return err
		}

		orderUpdateStatement, orderLineUpdateStatement, stockUpdateStatement := n.prepareStatements(newOrderID)

		if _, err := tx.Exec(stockUpdateStatement); err != nil {
			return fmt.Errorf("error in updating stock table: w=%d d=%d o=%d \n Err: %v", n.WarehouseID, n.DistrictID, newOrderID, err)
		}

		row := tx.QueryRow(orderUpdateStatement)
		if err := row.Scan(&result.OrderTimestamp); err != nil {
			return fmt.Errorf("error in inserting new order row: w=%d d=%d o=%d \n Err: %v", n.WarehouseID, n.DistrictID, newOrderID, err)
		}

		if _, err := tx.Exec(orderLineUpdateStatement); err != nil {
			return fmt.Errorf("error in inserting new order line rows: w=%d d=%d o=%d \n Err: %v", n.WarehouseID, n.DistrictID, newOrderID, err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error occured in updating the order/order lines/item pairs table. Err: %v", err)
	}

	result.TotalOrderAmount = n.TotalAmount * (1.0 + result.DistrictTax + result.WarehouseTax) * (1.0 - result.Customer.Discount)
	result.UniqueItems = n.UniqueItems
	result.OrderLineItems = n.NewOrderLineItems
	result.Print()

	return
}

// Print prints the formatted output of the NewOrder Transaction
func (o *Output) Print() {
	var newOrderString strings.Builder

	newOrderString.WriteString(fmt.Sprintf("Customer Identifier => W_ID = %d, D_ID = %d, C_ID = %d \n", o.Customer.WarehouseID, o.Customer.DistrictID, o.Customer.CustomerID))
	newOrderString.WriteString(fmt.Sprintf("Customer Info => Last Name: %s , Credit: %s , Discount: %0.6f \n", o.Customer.LastName, o.Customer.Credit, o.Customer.Discount))
	newOrderString.WriteString(fmt.Sprintf("Order Details: O_ID = %d , O_ENTRY_D = %s \n", o.OrderID, o.OrderTimestamp))
	newOrderString.WriteString(fmt.Sprintf("Total Unique Items: %d \n", o.UniqueItems))
	newOrderString.WriteString(fmt.Sprintf("Total Amount: %.2f \n", o.TotalOrderAmount))

	newOrderString.WriteString(fmt.Sprintf(" # \t ID \t Name (Supplier, Qty, Amount, Stock) \n"))
	idx := 1
	for key, value := range o.OrderLineItems {
		newOrderString.WriteString(fmt.Sprintf(" %02d \t %d \t %s (%d, %d, %.2f, %d) \n",
			idx,
			key,
			value.Name,
			value.SupplierWarehouseID,
			value.Quantity,
			value.Price*float64(value.Quantity),
			value.FinalStock,
		))
		idx++
	}

	fmt.Println(newOrderString.String())
}
