package tables

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/cdbconn"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
)

// ExecuteSQL executes the SQL file passed in the path variable
func ExecuteSQL(c config.Configuration, sqlFilePath string) error {
	log.Printf("Executing the SQL File: %s", sqlFilePath)

	db, err := cdbconn.CreateConnection(c.HostNode)
	if err != nil {
		panic("load function couldn't create a connection to the server")
	}

	sqlFile, err := os.Open(sqlFilePath)
	if err != nil {
		log.Fatalf("Err: %v", err)
		return errors.New("error occurred. Please check the logs")
	}
	defer sqlFile.Close()

	byteValue, _ := ioutil.ReadAll(sqlFile)

	var finalQueryBuilder strings.Builder
	finalQueryBuilder.WriteString(string(byteValue))

	for _, value := range strings.Split(finalQueryBuilder.String(), ";") {
		if _, err = db.Exec(value); err != nil {
			log.Println(value)
			log.Fatalf("Err: %v", err)
			return errors.New("error occurred. Please check the logs")
		}
	}

	log.Printf("Completed the SQL File: %s", sqlFilePath)
	return nil
}
