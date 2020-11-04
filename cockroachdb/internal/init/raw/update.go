package rawtables

import (
	"log"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/tables"
)

// UpdateParent udpates parent tables
func UpdateParent(c config.Configuration) error {
	log.Println("Updating parent tables...")

	log.Println("\nExecuting the SQL: scripts/sql/raw/update.sql")
	if err := tables.ExecuteSQL(c, "scripts/sql/raw/update.sql"); err != nil {
		log.Fatalf("error occured while loading raw tables. Err: %v", err)
		return err
	}

	log.Println("Updated all the parent tables...")
	return nil
}

// UpdatePartitions updates partitions of the tables
func UpdatePartitions(c config.Configuration) error {
	log.Println("Updating partitions of a table...")

	log.Println("Executing the SQL: scripts/sql/raw/update-partitions.sql")
	if err := tables.ExecuteSQLForPartitions(c, 10, 10, "scripts/sql/raw/update-partitions.sql"); err != nil {
		log.Fatalf("error occured while loading partitions. Err: %v", err)
		return err
	}

	log.Println("Updated all the partitions of the tables...")
	return nil
}
