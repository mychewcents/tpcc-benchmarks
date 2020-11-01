package main

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func loadRawDataset(db *sql.DB, file string) error {
	initSQL, err := os.Open(file)
	if err != nil {
		log.Fatalf("Err: %v", err)
		return errors.New("error occurred. Please check the logs")
	}
	defer initSQL.Close()

	byteValue, _ := ioutil.ReadAll(initSQL)

	var finalQueryBuilder strings.Builder
	finalQueryBuilder.WriteString(string(byteValue))

	for _, value := range strings.Split(finalQueryBuilder.String(), ";") {
		log.Println(value)

		if _, err = db.Exec(value); err != nil {
			log.Fatalf("Err: %v", err)
			return errors.New("error occurred. Please check the logs")
		}
	}
	return nil
}
