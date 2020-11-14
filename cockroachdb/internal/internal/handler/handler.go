package handler

import (
	"bufio"
)

// NewTransactionController defines the interface to handle transactions
type NewTransactionController interface {
	HandleTransaction(scanner *bufio.Scanner, args []string) bool
}
