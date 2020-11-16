package dao

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dbdatamodel"
)

// CustomerItemsPairDao interface of the Customer's Item pairs
type CustomerItemsPairDao interface {
	Insert(warehouseID, districtID, uniqueItems, customerID int, orderLineItems map[int]*dbdatamodel.OrderLineItem) error
	GetItemPairsString(warehouseID, districtID, customerID int) (string, error)
	GetCustomerIDsWithSamePairs(warehouseID, districtID int, itemPairs string) ([]int, error)
	InsertInitialPairs(warehouseID, districtID int, customerOrderIDs map[int]int, orderItems map[int][]int) (err error)
}

type customerItemsPairDaoImpl struct {
	db *sql.DB
}

// CreateCustomerItemsPairDao creates a new CustomerItemsDao
func CreateCustomerItemsPairDao(db *sql.DB) CustomerItemsPairDao {
	return &customerItemsPairDaoImpl{db: db}
}

// Insert inserts new customer items pairs during the new order transaction
func (cip *customerItemsPairDaoImpl) Insert(warehouseID, districtID, customerID, uniqueItems int, orderLineItems map[int]*dbdatamodel.OrderLineItem) error {
	var orderItemCustomerPair strings.Builder
	orderItemCustomerPairTable := fmt.Sprintf("ORDER_ITEMS_CUSTOMERS_%d_%d", warehouseID, districtID)

	sortedOrderItems := make([]int, uniqueItems)

	idx := 0
	for key := range orderLineItems {
		sortedOrderItems[idx] = key
		idx++
	}

	sort.Ints(sortedOrderItems)

	for i := 0; i < len(sortedOrderItems)-1; i++ {
		for j := i + 1; j < len(sortedOrderItems); j++ {
			orderItemCustomerPair.WriteString(fmt.Sprintf("(%d, %d, %d, %d, %d),", warehouseID, districtID, customerID, sortedOrderItems[i], sortedOrderItems[j]))
		}
	}

	orderItemCustomerPairString := orderItemCustomerPair.String()
	orderItemCustomerPairString = orderItemCustomerPairString[:len(orderItemCustomerPairString)-1]

	sqlStatement := fmt.Sprintf("UPSERT INTO %s (IC_W_ID, IC_D_ID, IC_C_ID, IC_I_1_ID, IC_I_2_ID) VALUES %s", orderItemCustomerPairTable, orderItemCustomerPairString)

	if _, err := cip.db.Exec(sqlStatement); err != nil {
		return fmt.Errorf("error occured in the item pairs for customers. Err: %v", err)
	}

	return nil
}

// GetItemPairsString returns the pairs of a customer as a string
func (cip *customerItemsPairDaoImpl) GetItemPairsString(warehouseID, districtID, customerID int) (itemPairsString string, err error) {
	var itemPairs strings.Builder

	sqlStatement := fmt.Sprintf("SELECT IC_I_1_ID, IC_I_2_ID FROM ORDER_ITEMS_CUSTOMERS_%d_%d WHERE IC_C_ID = %d", warehouseID, districtID, customerID)

	rows, err := cip.db.Query(sqlStatement)
	if err != nil {
		return "", fmt.Errorf("error in fetching the order line item pairs. Err: %v", err)
	}
	defer rows.Close()

	var itemID1, itemID2 int
	for rows.Next() {
		err := rows.Scan(&itemID1, &itemID2)
		if err != nil {
			return "", fmt.Errorf("error occurred in scanning the order line item pair. Err: %v", err)
		}
		itemPairs.WriteString(fmt.Sprintf("(%d, %d),", itemID1, itemID2))
	}

	itemPairsString = itemPairs.String()

	if len(itemPairsString) == 0 {
		return "", nil
	}

	itemPairsString = itemPairsString[:len(itemPairsString)-1]

	return
}

// GetCustomerIDsWithSamePairs fetches the customer ids with the same pairs as passed
func (cip *customerItemsPairDaoImpl) GetCustomerIDsWithSamePairs(warehouseID, districtID int, itemPairs string) (customerIDs []int, err error) {
	sqlStatement := fmt.Sprintf("SELECT IC_C_ID FROM ORDER_ITEMS_CUSTOMERS_%d_%d WHERE (IC_I_1_ID, IC_I_2_ID) IN (%s)", warehouseID, districtID, itemPairs)

	rows, err := cip.db.Query(sqlStatement)
	if err == sql.ErrNoRows {
		return customerIDs, nil
	}
	if err != nil {
		return nil, err
	}

	var cID int

	for rows.Next() {
		err := rows.Scan(&cID)
		if err != nil {
			return nil, err
		}
		customerIDs = append(customerIDs, cID)
	}

	return
}

// InsertInitialPairs inserts the pairs for the initial load
func (cip *customerItemsPairDaoImpl) InsertInitialPairs(warehouseID, districtID int, customerOrderIDs map[int]int, orderItems map[int][]int) (err error) {
	var customerItemsPairBuilder strings.Builder

	for cID, oID := range customerOrderIDs {
		sort.Ints(orderItems[oID])

		for i := 0; i < len(orderItems[oID])-1; i++ {
			for j := i + 1; j < len(orderItems[oID]); j++ {
				customerItemsPairBuilder.WriteString(fmt.Sprintf("(%d, %d, %d, %d, %d),", warehouseID, districtID, cID, orderItems[oID][i], orderItems[oID][j]))
			}
		}
	}

	customerItemsPairString := customerItemsPairBuilder.String()
	customerItemsPairString = customerItemsPairString[:len(customerItemsPairString)-1]

	sqlStatement := fmt.Sprintf("UPSERT INTO ORDER_ITEMS_CUSTOMERS_%d_%d (IC_W_ID, IC_D_ID, IC_C_ID, IC_I_1_ID, IC_I_2_ID) VALUES %s",
		warehouseID, districtID, customerItemsPairString)

	if _, err := cip.db.Exec(sqlStatement); err != nil {
		return err
	}

	return nil
}
