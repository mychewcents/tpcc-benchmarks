package controller

import (
	"fmt"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/common"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/handler"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/internal/model"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/internal/service"
	"strconv"
)

type StockLevelController interface {
	handler.TransactionHandler
}

type stockLevelControllerImpl struct {
	s service.StockLevelService
}

func NewStockLevelController(cassandraSession *common.CassandraSession) StockLevelController {
	return &stockLevelControllerImpl{
		s: service.NewStockLevelService(cassandraSession),
	}
}

func (s *stockLevelControllerImpl) HandleTransaction(cmd []string) {
	request := makeStockLevelRequest(cmd)
	response, _ := s.s.ProcessStockLevelTransaction(request)
	printStockLevelResponse(response)
}

func makeStockLevelRequest(cmd []string) *model.StockLevelRequest {
	wId, _ := strconv.Atoi(cmd[1])
	dId, _ := strconv.Atoi(cmd[2])
	t, _ := strconv.Atoi(cmd[3])
	l, _ := strconv.Atoi(cmd[4])

	return &model.StockLevelRequest{
		WId:            wId,
		DId:            dId,
		Threshold:      t,
		NoOfLastOrders: l,
	}
}

func printStockLevelResponse(r *model.StockLevelResponse) {
	fmt.Println("*********************** Stock Level Transaction Output ***********************")
	fmt.Println("Total number of items in examined orders where its stock quantity at warehouse is below the threshold: ", r.Count)
	fmt.Println()
}

func (s *stockLevelControllerImpl) Close() error {
	panic("implement me")
}
