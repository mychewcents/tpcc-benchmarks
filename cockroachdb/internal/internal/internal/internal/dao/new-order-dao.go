package dao

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
)

// NewOrderDao provides the interface for the required functions
type NewOrderDao interface {
	GetNewOrderIDAndTaxRates(req *models.NewOrder) (newOrderID int, wTax, dTax float64, err error)
	GetCustomerInformation(req *models.NewOrder) (cLastName, cCredit string, cDiscount float64, err error)
	InsertOrderPairItems(req *models.NewOrder) error
	GetItemDetails(tx *sql.Tx, req *models.NewOrder) error
	PrepareStatements(orderID int, req *models.NewOrder) (orderUpdateStatement, orderLineUpdateStatement, stockUpdateStatement string)
}

// NewOrderDaoImpl dao implementation
type NewOrderDaoImpl struct {
	db *sql.DB
}

// GetNewNewOrderDao gets the new DAO implementation for the NewOrder
func GetNewNewOrderDao(db *sql.DB) NewOrderDao {
	return &NewOrderDaoImpl{db: db}
}

// GetNewOrderIDAndTaxRates gets the new order id and the tax rates
func (n *NewOrderDaoImpl) GetNewOrderIDAndTaxRates(req *models.NewOrder) (newOrderID int, wTax, dTax float64, err error) {
	sqlStatement := fmt.Sprintf("UPDATE District SET D_NEXT_O_ID = D_NEXT_O_ID + 1 WHERE D_W_ID = $1 AND D_ID = $2 RETURNING D_NEXT_O_ID, D_TAX, D_W_TAX")

	row := n.db.QueryRow(sqlStatement, req.WarehouseID, req.DistrictID)
	if err := row.Scan(&newOrderID, &dTax, &wTax); err != nil {
		return 0, 0.0, 0.0, fmt.Errorf("error occured in updating the district table for the next order id. Err: %v", err)
	}

	return
}

// GetCustomerInformation gets the customer information
func (n *NewOrderDaoImpl) GetCustomerInformation(req *models.NewOrder) (cLastName, cCredit string, cDiscount float64, err error) {
	sqlStatement := fmt.Sprintf("SELECT C_LAST, C_CREDIT, C_DISCOUNT FROM CUSTOMER WHERE C_W_ID = $1 AND C_D_ID = $2 AND C_ID = $3")

	row := n.db.QueryRow(sqlStatement, req.WarehouseID, req.DistrictID, req.CustomerID)
	if err := row.Scan(&cLastName, &cCredit, &cDiscount); err != nil {
		return "", "", 0.0, fmt.Errorf("error occured in getting the customer details. Err: %v", err)
	}

	return
}

// InsertOrderPairItems inserts new order item pairs
func (n *NewOrderDaoImpl) InsertOrderPairItems(req *models.NewOrder) error {
	var orderItemCustomerPair strings.Builder
	orderItemCustomerPairTable := fmt.Sprintf("ORDER_ITEMS_CUSTOMERS_%d_%d", req.WarehouseID, req.DistrictID)

	sortedOrderItems := make([]int, req.UniqueItems)

	idx := 0
	for key := range req.NewOrderLineItems {
		sortedOrderItems[idx] = key
		idx++
	}

	sort.Ints(sortedOrderItems)

	for i := 0; i < len(sortedOrderItems)-1; i++ {
		for j := i + 1; j < len(sortedOrderItems); j++ {
			orderItemCustomerPair.WriteString(fmt.Sprintf("(%d, %d, %d, %d, %d),", req.WarehouseID, req.DistrictID, req.CustomerID, sortedOrderItems[i], sortedOrderItems[j]))
		}
	}

	sqlStatement := fmt.Sprintf("UPSERT INTO %s (IC_W_ID, IC_D_ID, IC_C_ID, IC_I_1_ID, IC_I_2_ID) VALUES %s", orderItemCustomerPairTable, orderItemCustomerPair)
	sqlStatement = sqlStatement[0 : len(sqlStatement)-1]

	if _, err := n.db.Exec(sqlStatement); err != nil {
		return fmt.Errorf("error occured in the item pairs for customers. Err: %v", err)
	}

	return nil
}

// GetItemDetails gets the items details
func (n *NewOrderDaoImpl) GetItemDetails(tx *sql.Tx, req *models.NewOrder) error {

	var itemsWhereClause strings.Builder

	for key, value := range req.NewOrderLineItems {
		itemsWhereClause.WriteString(fmt.Sprintf("(%d, %d),", value.SupplierWarehouseID, key))
	}

	itemsWhereClauseString := itemsWhereClause.String()
	itemsWhereClauseString = itemsWhereClauseString[:len(itemsWhereClauseString)-1]

	sqlStatement := fmt.Sprintf("SELECT S_I_ID, S_I_NAME, S_I_PRICE, S_QUANTITY, S_YTD, S_ORDER_CNT, S_DIST_%02d FROM STOCK WHERE (S_W_ID, S_I_ID) IN %s",
		req.DistrictID, itemsWhereClauseString)
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

		if value, ok := req.NewOrderLineItems[id]; ok {
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
			req.TotalAmount += value.Amount
		}
	}

	return nil
}

// PrepareStatements prepares the statements for the DB update
func (n *NewOrderDaoImpl) PrepareStatements(orderID int, req *models.NewOrder) (orderUpdateStatement, orderLineUpdateStatement, stockUpdateStatement string) {
	var orderLineEntries, stockOrderItemIdentifiers, stockQuantityUpdates, stockYTDUpdates, stockOrderCountUpdates, stockRemoteCountUpdates strings.Builder

	var itemIdentifier string
	whenClauseFormat := "WHEN %s THEN %d "

	idx := 0
	for key, value := range req.NewOrderLineItems {
		orderLineEntries.WriteString(
			fmt.Sprintf("(%d, %d, %d, %d, %d, %d, %d, %0.2f, '%s'),",
				orderID,
				req.DistrictID,
				req.WarehouseID,
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
		req.WarehouseID,
		req.DistrictID,
		orderID,
		req.WarehouseID,
		req.DistrictID,
		req.CustomerID,
		req.UniqueItems,
		req.IsOrderLocal,
		req.TotalAmount,
	)

	orderLineUpdateStatement = fmt.Sprintf("INSERT INTO ORDER_LINE_%d_%d (OL_O_ID, OL_D_ID, OL_W_ID, OL_NUMBER, OL_I_ID, OL_SUPPLY_W_ID, OL_QUANTITY, OL_AMOUNT, OL_DIST_INFO) VALUES %s",
		req.WarehouseID, req.DistrictID, orderLineEntriesString)

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
