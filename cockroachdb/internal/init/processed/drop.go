package processedtables

import (
	"log"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/tables"
)

// DropParent dropes parent tables
func DropParent(c config.Configuration) error {
	log.Println("Dropping parent tables...")

	log.Println("\nExecuting the SQL: scripts/sql/processed/drop.sql")
	if err := tables.ExecuteSQL(c, "scripts/sql/processed/drop.sql"); err != nil {
		log.Fatalf("error occured while dropping processed tables. Err: %v", err)
		return err
	}

	log.Println("Dropped all the parent tables...")
	return nil
}

// DropPartitions dropes partitioned tables
func DropPartitions(c config.Configuration) error {
	log.Println("Dropping partitions of a table...")

	log.Println("Executing the SQL: scripts/sql/processed/drop-partitions.sql")
	if err := tables.ExecuteSQLForPartitions(c, 10, 10, "scripts/sql/processed/drop-partitions.sql"); err != nil {
		log.Fatalf("error occured while dropping partitions. Err: %v", err)
		return err
	}

	log.Println("Dropped all the partitions of the tables...")
	return nil
}
