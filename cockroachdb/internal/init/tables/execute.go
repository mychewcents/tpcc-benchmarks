package tables

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// ExecuteSQL executes the SQL file passed in the path variable
func ExecuteSQL(db *sql.DB, sqlFilePath string) error {
	log.Printf("Executing the SQL File: %s", sqlFilePath)

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
			log.Fatalln(value)
			log.Fatalf("Err: %v", err)
			return errors.New("error occurred. Please check the logs")
		}
	}

	log.Printf("Completed the SQL File: %s", sqlFilePath)
	return nil
}
