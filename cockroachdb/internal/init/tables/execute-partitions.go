package tables

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// ExecuteSQLForPartitions executes the SQL for the required partitions
func ExecuteSQLForPartitions(db *sql.DB, warehouses, districts int, sqlFilePath string) error {
	log.Printf("Executing the SQL File: %s", sqlFilePath)

	sqlFile, err := os.Open(sqlFilePath)
	if err != nil {
		log.Fatalf("Err: %v", err)
		return errors.New("error occurred. Please check the logs")
	}
	defer sqlFile.Close()

	byteValue, _ := ioutil.ReadAll(sqlFile)
	baseSQLStatement := string(byteValue)

	errFound := false
	for i := 1; i <= warehouses; i++ {
		for j := 1; j <= districts; j++ {
			finalSQLStatement := strings.ReplaceAll(baseSQLStatement, "WID", strconv.Itoa(i))
			finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "DID", strconv.Itoa(j))

			if _, err := db.Exec(finalSQLStatement); err != nil {
				log.Fatalf("Query: %s", finalSQLStatement)
				log.Fatalf("Err: %v", err)
				errFound = true
			}
		}
	}

	if errFound {
		return errors.New("error was found. Please check the logs")
	}

	log.Printf("Completed the SQL File: %s", sqlFilePath)
	return nil
}
