package processedtables

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/cdbconn"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/tables"
)

// LoadParent loads parent tables
func LoadParent(c config.Configuration) error {
	log.Println("Loading parent tables...")

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

	ch := make(chan bool, 10)
	for w := 1; w <= 10; w++ {
		go loadPartitionsParallel(c, w, baseSQLStatement, ch)

	}

	for i := 0; i < 10; i++ {
		<-ch
	}

	log.Println("Loaded all the partitions of the tables...")
	return nil
}

func loadPartitionsParallel(c config.Configuration, w int, baseSQLStatement string, ch chan bool) {
	db, err := cdbconn.CreateConnection(c.HostNode)
	if err != nil {
		panic("load function couldn't create a connection to the server")
	}

	for d := 1; d <= 10; d++ {
		finalSQLStatement := strings.ReplaceAll(baseSQLStatement, "ORDERS_FILE_PATH", fmt.Sprintf("order/%d_%d", w, d))
		finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "ORDER_LINE_FILE_PATH", fmt.Sprintf("orderline/%d_%d", w, d))
		finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "ORDER_ITEMS_CUSTOMERS_FILE_PATH", fmt.Sprintf("itempairs/%d_%d", w, d))
		finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "WID", strconv.Itoa(w))
		finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "DID", strconv.Itoa(d))

		_, err := db.Exec(finalSQLStatement)
		if err != nil {
			log.Fatalf("couldn't load the table: %d %d. Err: %v", w, d, err)
			ch <- false
		}
	}

	log.Printf("Completed Partition: %d", w)
	ch <- true
}
