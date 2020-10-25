package caller

import (
	"bufio"
	"database/sql"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/neworder"
)

// ProcessRequest Calls the required DB function
func ProcessRequest(db *sql.DB, scanner *bufio.Scanner, transactionArgs []string) error {
	// var err error

	switch transactionArgs[0] {
	case "N":
		neworder.ProcessTransaction(db, scanner, transactionArgs[1:])
	case "P":
	case "D":
	case "O":
	case "S":
	case "I":
	case "T":
	case "R":

	}
	return nil
}
