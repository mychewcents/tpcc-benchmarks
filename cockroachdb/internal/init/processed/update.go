package processedtables

import (
	"log"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/connection/config"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/tables"
)

// UpdateParent udpates parent tables
func UpdateParent(c config.Configuration) error {
	log.Printf("No update required")
	return nil
}

// UpdatePartitions updates partitions of the tables
func UpdatePartitions(c config.Configuration) error {
	log.Println("Updating partitions of tables...")

	if err := tables.ExecuteSQLForPartitions(c, 10, 10, "scripts/sql/processed/update-partitions.sql"); err != nil {
		log.Fatalf("error occured while updating partitions. Err: %v", err)
		return err
	}

	log.Println("Updated all the partitions of the tables...")
	return nil
}
