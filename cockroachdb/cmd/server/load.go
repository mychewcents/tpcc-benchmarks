package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/cdbconn"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/init/tables"
)

func load(c config.Configuration) {

	db, err := cdbconn.CreateConnection(c.HostNode)
	if err != nil {
		panic("load function couldn't create a connection to the server")
	}

	fmt.Printf("Executing the SQL: scripts/sql/drop-partitions.sql")
	if err := tables.ExecuteSQLForPartitions(db, 10, 10, "scripts/sql/drop-partitions.sql"); err != nil {
		fmt.Println(err)
		return
	}

	sqlScripts := []string{
		"scripts/sql/drop-raw.sql",
		"scripts/sql/create-raw.sql",
		"scripts/sql/load-raw.sql",
		"scripts/sql/update-raw.sql",
	}

	for _, value := range sqlScripts {
		fmt.Printf("\nExecuting the SQL: %s", value)
		if err := tables.ExecuteSQL(db, value); err != nil {
			fmt.Println(err)
			return
		}
	}

	sqlScripts = []string{
		"scripts/sql/create-partitions.sql",
		"scripts/sql/load-partitions.sql",
		"scripts/sql/update-partitions.sql",
	}

	for _, value := range sqlScripts {
		fmt.Printf("\nExecuting the SQL: %s", value)
		if err := tables.ExecuteSQLForPartitions(db, 10, 10, value); err != nil {
			fmt.Println(err)
			return
		}
	}

	if err := loadOrderItemsCustomerPair(db, 10); err != nil {
		log.Fatalf("error in loadOrderItemsCustomerPair. Err: %v", err)
	}
	log.Println("Initialization Complete!")
	fmt.Println("\nInitialization Complete!")
}

func loadOrderItemsCustomerPair(db *sql.DB, warehouses int) error {

	return nil
}
