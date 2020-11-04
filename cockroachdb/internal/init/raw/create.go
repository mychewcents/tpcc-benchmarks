package rawtables

import (
	"log"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/tables"
)

// CreateParent creates tables
func CreateParent(c config.Configuration) error {
	log.Println("Creating tables...")

	log.Println("\nExecuting the SQL: scripts/sql/raw/create.sql")
	if err := tables.ExecuteSQL(c, "scripts/sql/raw/create.sql"); err != nil {
		log.Fatalf("error occured while creating raw tables. Err: %v", err)
		return err
	}

	log.Println("Created all the tables...")
	return nil
}

// CreatePartitions creates partitioned tables
func CreatePartitions(c config.Configuration) error {
	log.Println("Creating partitions of tables...")

	log.Println("Executing the SQL: scripts/sql/raw/create-partitions.sql")
	if err := tables.ExecuteSQLForPartitions(c, 10, 10, "scripts/sql/raw/create-partitions.sql"); err != nil {
		log.Fatalf("error occured while creating partitions. Err: %v", err)
		return err
	}

	log.Println("Created all the partitions of the tables...")
	return nil
}
