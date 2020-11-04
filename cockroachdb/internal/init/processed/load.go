package processedtables

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/cdbconn"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/tables"
)

// LoadParent loads parent tables
func LoadParent(c config.Configuration) error {
	log.Println("Loading parent tables...")

	log.Println("\nExecuting the SQL: scripts/sql/processed/load.sql")
	if err := tables.ExecuteSQL(c, "scripts/sql/processed/load.sql"); err != nil {
		log.Fatalf("error occured while loading processed tables. Err: %v", err)
		return err
	}

	log.Println("Loaded all the parent tables...")
	return nil
}

// LoadPartitions loads partitioned tables
func LoadPartitions(c config.Configuration) error {
	log.Println("Loading partitions of a table...")

	log.Println("Executing the SQL: scripts/sql/processed/load-partitions.sql")
	sqlScript := "scripts/sql/processed/load-partitions.sql"

	sqlFile, err := os.Open(sqlScript)
	if err != nil {
		log.Fatalf("Err: %v", err)
		return err
	}
	defer sqlFile.Close()

	byteValue, _ := ioutil.ReadAll(sqlFile)
	baseSQLStatement := string(byteValue)

	db, err := cdbconn.CreateConnection(c.HostNode)
	if err != nil {
		panic("load function couldn't create a connection to the server")
	}

	for w := 1; w <= 10; w++ {
		for d := 1; d <= 10; d++ {
			finalSQLStatement := strings.ReplaceAll(baseSQLStatement, "ORDERS_FILE_PATH", fmt.Sprintf("order/%d_%d", w, d))
			finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "ORDER_LINE_FILE_PATH", fmt.Sprintf("orderline/%d_%d", w, d))
			finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "ORDER_ITEMS_CUSTOMERS_FILE_PATH", fmt.Sprintf("itempairs/%d_%d", w, d))

			_, err := db.Exec(finalSQLStatement)
			if err != nil {
				log.Fatalf("couldn't load the table: %d %d. Err: %v", w, d, err)
			}
			log.Printf("Completed Partition: %d %d", w, d)
		}
	}

	log.Println("Loaded all the partitions of the tables...")
	return nil
}
