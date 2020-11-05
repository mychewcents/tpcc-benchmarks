package caller

import (
	"bufio"
	"database/sql"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/transactions/delivery"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/transactions/neworder"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/transactions/orderstatus"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/transactions/payment"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/transactions/popularitem"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/transactions/relatedcustomer"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/transactions/stocklevel"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/transactions/topbalance"
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
		return relatedcustomer.ProcessTransaction(db, nil, transactionArgs[1:])
	}
	return false
}
