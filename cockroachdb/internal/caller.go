package caller

import (
	"bufio"
	"database/sql"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/neworder"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/popularitem"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/stocklevel"
)

// ProcessRequest Calls the required DB function
func ProcessRequest(db *sql.DB, scanner *bufio.Scanner, transactionArgs []string) {

	switch transactionArgs[0] {
	case "N":
		neworder.ProcessTransaction(db, scanner, transactionArgs[1:])
	case "P":
	case "D":
	case "O":
	case "S":
		stocklevel.ProcessTransaction(db, nil, transactionArgs[1:])
	case "I":
		popularitem.ProcessTransaction(db, nil, transactionArgs[1:])
	case "T":
	case "R":

	}

}
