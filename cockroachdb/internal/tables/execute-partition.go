package tables

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/cdbconn"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
)

// ExecuteSQLForPartitions executes the SQL for the required partitions
func ExecuteSQLForPartitions(c config.Configuration, warehouses, districts int, sqlFilePath string) error {
	log.Printf("Executing the SQL File: %s\n", sqlFilePath)

	sqlFile, err := os.Open(sqlFilePath)
	if err != nil {
		log.Fatalf("Err: %v", err)
		return errors.New("error occurred. Please check the logs")
	}
	defer sqlFile.Close()

	byteValue, _ := ioutil.ReadAll(sqlFile)
	baseSQLStatement := string(byteValue)

	ch := make(chan bool, 100)
	errFound := false
	for w := 1; w <= warehouses; w++ {
		for d := 1; d <= districts; d++ {
			finalSQLStatement := baseSQLStatement
			finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "WID", strconv.Itoa(w))
			finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "DID", strconv.Itoa(d))

			go executeParallel(c, w, d, finalSQLStatement, ch)
		}
	}

	for i := 0; i < 100; i++ {
		<-ch
	}

	if errFound {
		return errors.New("error was found. Please check the logs")
	}

	log.Printf("Completed the SQL File: %s", sqlFilePath)
	return nil
}

func executeParallel(c config.Configuration, w, d int, finalSQLStatement string, ch chan bool) {
	db, err := cdbconn.CreateConnection(c.HostNode)
	if err != nil {
		panic("load function couldn't create a connection to the server")
	}

	if _, err := db.Exec(finalSQLStatement); err != nil {
		log.Println(finalSQLStatement)
		log.Fatalf("Err: %v", err)
		ch <- false
	}
	log.Printf("Executed partition: %d %d", w, d)
	ch <- true
}
