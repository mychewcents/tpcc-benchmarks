package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/cdbconn"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/init/tables"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/logging"
)

var db *sql.DB

func init() {
	var err error

	if len(os.Args) != 2 {
		panic("Missing configuration file path")
	}
	db, err = cdbconn.CreateConnection(os.Args[1])
	if err != nil {
		panic(err)
	}

	if err := logging.SetupLogOutput("init", "logs"); err != nil {
		panic(err)
	}
}

func main() {
	sqlScripts := []string{
		"scripts/sql/drop-raw.sql",
		"scripts/sql/create-raw.sql",
		"scripts/sql/load-raw.sql",
		"scripts/sql/update-raw.sql",
	}

	for _, value := range sqlScripts {
		if err := tables.ExecuteSQL(db, value); err != nil {
			fmt.Println(err)
			return
		}
	}

	sqlScripts = []string{
		"scripts/sql/drop-partitions.sql",
		"scripts/sql/create-partitions.sql",
		"scripts/sql/load-partitions.sql",
		"scripts/sql/update-partitions.sql",
	}

	for _, value := range sqlScripts {
		if err := tables.ExecuteSQLForPartitions(db, 10, 10, value); err != nil {
			fmt.Println(err)
			return
		}
	}

	log.Println("Initialization Complete!")
	fmt.Println("Initialization Complete!")
}
