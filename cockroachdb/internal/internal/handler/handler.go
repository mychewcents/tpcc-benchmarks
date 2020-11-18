package handler

import (
	"bufio"
)

// NewTransactionController defines the interface to handle transactions
type NewTransactionController interface {
	HandleTransaction(scanner *bufio.Scanner, args []string) bool
}

// NewTablesController defines the interface to the inital table load
type NewTablesController interface {
	LoadProcessedTables() error
	LoadRawTables() error
	ExportTables() error
}
