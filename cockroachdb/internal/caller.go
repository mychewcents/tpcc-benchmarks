package caller

import (
	"bufio"
	"database/sql"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/orderstatus"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/payment"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/neworder"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/popularitem"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/stocklevel"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/topbalance"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/delivery"
)

// ProcessRequest Calls the required DB function
func ProcessRequest(db *sql.DB, scanner *bufio.Scanner, transactionArgs []string) bool {
	switch transactionArgs[0] {
	case "N":
		return neworder.ProcessTransaction(db, scanner, transactionArgs[1:])
	case "P":
		return payment.ProcessTransaction(db, nil, transactionArgs[1:])
	case "D":
		return delivery.ProcessTransaction(db, nil, transactionArgs[1:])
	case "O":
		return orderstatus.ProcessTransaction(db, nil, transactionArgs[1:])
	case "S":
		return stocklevel.ProcessTransaction(db, nil, transactionArgs[1:])
	case "I":
		return popularitem.ProcessTransaction(db, nil, transactionArgs[1:])
	case "T":
		return topbalance.ProcessTransaction(db, nil)
	case "R":

	}
	return false
}
