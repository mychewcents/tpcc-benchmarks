package handler

import (
	"io"
)

type TransactionHandler interface {
	HandleTransaction(cmd []string)
	io.Closer
}
