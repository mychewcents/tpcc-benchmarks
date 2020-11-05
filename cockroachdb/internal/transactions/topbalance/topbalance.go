package topbalance

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/cockroachdb/cockroach-go/crdb"
)

type txnOutput struct {
	first, middle, last, warehouseName, districtName string
	warehouseID, districtID                          int
	balance                                          float64
}

// ProcessTransaction gets the top balance customer
func ProcessTransaction(db *sql.DB, scanner *bufio.Scanner) bool {
	log.Printf("Starting the Top Balance Transaction")

	if err := execute(db); err != nil {
		log.Fatalf("error occurred while executing the top balance transaction. Err: %v", err)
		return false
	}

	log.Printf("Completed the Top Balance Transaction")
	return true
}

func execute(db *sql.DB) error {
	log.Printf("Executing the transaction with the input data...")

	customerQuery := "SELECT C_FIRST, C_MIDDLE, C_LAST, C_W_ID, C_D_ID, C_BALANCE FROM CUSTOMER ORDER BY C_BALANCE DESC LIMIT 10"
	districtQuery := "SELECT D_NAME FROM DISTRICT WHERE D_W_ID=%d AND D_ID=%d"
	warehouseQuery := "SELECT W_NAME FROM WAREHOUSE WHERE W_ID=%d"

	var customers []txnOutput
	err := crdb.ExecuteTx(context.Background(), db, nil, func(tx *sql.Tx) error {
		rows, err := tx.Query(customerQuery)
		if err != nil {
			return fmt.Errorf("error occured while reading the customer details. Err: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			var customer txnOutput
			if err := rows.Scan(&customer.first, &customer.middle, &customer.last, &customer.warehouseID, &customer.districtID, &customer.balance); err != nil {
				return fmt.Errorf("error occured while scanning the customer details. Err: %v", err)
			}
			customers = append(customers, customer)
		}
		for i := range customers {
			if err := db.QueryRow(fmt.Sprintf(districtQuery, customers[i].warehouseID, customers[i].districtID)).Scan(&customers[i].districtName); err != nil {
				return fmt.Errorf("error occured while reading the district name. Err: %v", err)
			}
			if err := db.QueryRow(fmt.Sprintf(warehouseQuery, customers[i].warehouseID)).Scan(&customers[i].warehouseName); err != nil {
				return fmt.Errorf("error occured while reading the warehouse name. Err: %v", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error occured while updating the tables. Err: %v", err)
	}

	// for _, customer := range customers {
	// 	fmt.Println(fmt.Sprintf("Customer: %s %s %s, Balance: %f, Warehouse: %s, District: %s",
	// 		customer.first, customer.middle, customer.last, customer.balance, customer.warehouseName, customer.districtName))
	// }

	log.Printf("Executing the transaction with the input data...")
	return nil
}
