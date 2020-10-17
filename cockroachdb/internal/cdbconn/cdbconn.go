package cdbconn

import (
	"database/sql"
	"log"
)

// CreateConnection returns a DB connection object
func CreateConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgresql://root@localhost:30000/default?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB unreachable:", err)
	}
	return db, err
}
