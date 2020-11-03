package cdbconn

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"

	// Because we're using the Postgres Driver
	_ "github.com/lib/pq"
)

// CreateConnection returns a DB connection object
func CreateConnection(server config.Server) (*sql.DB, error) {
	connString := fmt.Sprintf("postgres://%s@%s:%d/%s?sslmode=disable",
		server.Username, server.Host, server.Port, server.Database)

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
