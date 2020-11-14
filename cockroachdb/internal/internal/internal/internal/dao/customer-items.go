package dao

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
)

// CustomerItemsPairDao interface of the Customer's Item pairs
type CustomerItemsPairDao interface {
	Insert(warehouseID, districtID, uniqueItems, customerID int, orderLineItems map[int]*models.NewOrderOrderLineItem) error
}

type customerItemsPairDaoImpl struct {
	db *sql.DB
}

// CreateCustomerItemsPairDao creates a new CustomerItemsDao
func CreateCustomerItemsPairDao(db *sql.DB) CustomerItemsPairDao {
	return &customerItemsPairDaoImpl{db: db}
}

func (cip *customerItemsPairDaoImpl) Insert(warehouseID, districtID, customerID, uniqueItems int, orderLineItems map[int]*models.NewOrderOrderLineItem) error {
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
