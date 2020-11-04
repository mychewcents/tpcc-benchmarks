package rawtables

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/cdbconn"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/tables"
)

// LoadParent loads parent tables
func LoadParent(c config.Configuration) error {
	log.Println("Loading parent tables...")

	log.Println("\nExecuting the SQL: scripts/sql/raw/load.sql")
	if err := tables.ExecuteSQL(c, "scripts/sql/raw/load.sql"); err != nil {
		log.Fatalf("error occured while loading raw tables. Err: %v", err)
		return err
	}

	log.Println("Loaded all the parent tables...")
	return nil
}

// LoadPartitions loads partitioned tables
func LoadPartitions(c config.Configuration) error {
	log.Println("Loading partitions of a table...")

	log.Println("Executing the SQL: scripts/sql/raw/load-partitions.sql")
	if err := tables.ExecuteSQLForPartitions(c, 10, 10, "scripts/sql/raw/load-partitions.sql"); err != nil {
		log.Fatalf("error occured while loading partitions. Err: %v", err)
		return err
	}

	if err := loadOrderItemsCustomerPair(c); err != nil {
		log.Fatalf("error in loadOrderItemsCustomerPair. Err: %v", err)
	}

	log.Println("Loaded all the partitions of the tables...")
	return nil
}

func loadOrderItemsCustomerPair(c config.Configuration) error {
	log.Println("Executing the load of Items Customer Pair")
	ch := make(chan bool, 100)

	for w := 1; w <= 10; w++ {
		for d := 1; d <= 10; d++ {
			orderItemCustomerPairTable := fmt.Sprintf("ORDER_ITEMS_CUSTOMERS_%d_%d", w, d)
			orderLineTable := fmt.Sprintf("ORDER_LINE_%d_%d", w, d)
			orderTable := fmt.Sprintf("ORDERS_%d_%d", w, d)

			go loadOrderItemsCustomerPairParallel(c, w, d, orderTable, orderLineTable, orderItemCustomerPairTable, ch)
		}
	}

	for i := 1; i <= 100; i++ {
		<-ch
	}

	log.Println("Completed the load of Items Customer Pair")

	return nil
}

func loadOrderItemsCustomerPairParallel(c config.Configuration, w, d int, orderTable, orderLineTable, orderItemCustomerPairTable string, ch chan bool) {
	db, err := cdbconn.CreateConnection(c.HostNode)
	if err != nil {
		panic("load function couldn't create a connection to the server")
	}

	for c := 1; c <= 3000; c++ {
		var orderID, orderLineItemsCount int

		sqlStatement := fmt.Sprintf("SELECT O_ID, O_OL_CNT FROM %s WHERE O_C_ID = %d", orderTable, c)

		row := db.QueryRow(sqlStatement)
		if err := row.Scan(&orderID, &orderLineItemsCount); err != nil {
			log.Fatalf("error in getting the order id for w = %d d = %d c = %d", w, d, c)
		}

		sqlStatement = fmt.Sprintf("SELECT OL_I_ID FROM %s WHERE OL_O_ID = %d", orderLineTable, orderID)
		rows, err := db.Query(sqlStatement)
		if err != nil {
			log.Fatalf("error in getting the order line items for w = %d d = %d orderid = %d", w, d, orderID)
		}

		var orderLineItem, totalUniqueItems int
		orderLineItems := make([]int, orderLineItemsCount)
		orderLineItemMap := make(map[int]bool)
		for rows.Next() {
			err := rows.Scan(&orderLineItem)
			if err != nil {
				log.Fatalf("error reading the output of the order line number for w = %d d = %d orderid = %d", w, d, orderID)
			}

			if _, ok := orderLineItemMap[orderLineItem]; !ok {
				orderLineItemMap[orderLineItem] = true
				orderLineItems[totalUniqueItems] = orderLineItem
				totalUniqueItems++
			}
		}
		orderLineItems = orderLineItems[:totalUniqueItems]
		sort.Ints(orderLineItems)

		var orderItemCustomerPair strings.Builder

		for i := 0; i < len(orderLineItems)-1; i++ {
			for j := i + 1; j < len(orderLineItems); j++ {
				orderItemCustomerPair.WriteString(fmt.Sprintf("(%d, %d, %d, %d, %d),", w, d, c, orderLineItems[i], orderLineItems[j]))
			}
		}

		sqlStatement = fmt.Sprintf("UPSERT INTO %s (IC_W_ID, IC_D_ID, IC_C_ID, IC_I_1_ID, IC_I_2_ID) VALUES %s", orderItemCustomerPairTable, orderItemCustomerPair.String())
		sqlStatement = sqlStatement[0 : len(sqlStatement)-1]

		if _, err := db.Exec(sqlStatement); err != nil {
			fmt.Println(err)
		}
	}

	log.Printf("Executed partition: %d %d", w, d)
	ch <- true
}
