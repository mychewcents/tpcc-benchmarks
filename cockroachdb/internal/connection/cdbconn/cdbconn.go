package cdbconn

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	// Because we're using the Postgres Driver
	_ "github.com/lib/pq"
)

type configuration struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Database string `json:"database"`
}

// CreateConnection returns a DB connection object
func CreateConnection(path string) (*sql.DB, error) {

	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	byteValue, _ := ioutil.ReadAll(configFile)

	var config configuration

	if err = json.Unmarshal(byteValue, &config); err != nil {
		return nil, err
	}

	configFile.Close()

	connString := fmt.Sprintf("postgres://%s@%s:%d/%s?sslmode=disable",
		config.Username, config.Host, config.Port, config.Database)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB unreachable:", err)
		return nil, err
	}
	return db, err
}
