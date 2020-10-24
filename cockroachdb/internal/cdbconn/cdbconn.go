package cdbconn

import (
	"database/sql"
	"fmt"
	"log"

	// Because we're using the Postgres Driver
	_ "github.com/lib/pq"
)

// CreateConnection returns a DB connection object
func CreateConnection(connAddr, port, database, username string) (*sql.DB, error) {

	connString := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", username, connAddr, port, database)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB unreachable:", err)
	}
	return db, err
}
