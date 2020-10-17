package router

import (
	"bufio"
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/controller"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/handler"
	"strings"
)

type TransactionRouter interface {
	HandleCommand(command string)
}

type transactionRouterImpl struct {
	handlers map[string]handler.TransactionHandler
}

func NewTransactionRouter(cluster *gocql.ClusterConfig, reader *bufio.Reader) TransactionRouter {
	router := transactionRouterImpl{
		handlers: make(map[string]handler.TransactionHandler, 0),
	}
	router.registerHandlers(cluster, reader)
	return &router
}

func (t *transactionRouterImpl) HandleCommand(command string) {
	commandSplit := strings.Split(strings.Trim(command, "\n"), ",")
	t.handlers[commandSplit[0]].HandleTransaction(commandSplit)
}

func (t *transactionRouterImpl) registerHandlers(cluster *gocql.ClusterConfig, reader *bufio.Reader) {
	t.handlers["N"] = controller.NewNewOrderTransactionController(cluster, reader)
	t.handlers["P"] = controller.NewPaymentController(cluster, reader)
	t.handlers["D"] = controller.NewDeliveryTransactionController(cluster, reader)
	t.handlers["O"] = controller.NewOrderStatusTransactionController(cluster, reader)
	t.handlers["S"] = controller.NewStockLevelController(cluster, reader)
	t.handlers["I"] = controller.NewPopularItemController(cluster, reader)
	t.handlers["T"] = controller.NewTopBalanceController(cluster, reader)
	t.handlers["R"] = controller.NewRelatedCustomerController(cluster, reader)
}
