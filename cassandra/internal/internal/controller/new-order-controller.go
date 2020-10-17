package controller

import (
	"bufio"
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/handler"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/service"
	"strconv"
	"strings"
)

type NewOrderController interface {
	handler.TransactionHandler
}

type newOrderControllerImpl struct {
	s service.NewOrderService
	r *bufio.Reader
}

func NewNewOrderTransactionController(cluster *gocql.ClusterConfig, reader *bufio.Reader) NewOrderController {
	return &newOrderControllerImpl{
		s: service.NewNewOrderService(cluster),
		r: reader,
	}
}

func (n *newOrderControllerImpl) HandleTransaction(cmd []string) {
	request := makeNewOrderRequest(cmd, n.r)
	n.s.ProcessNewOrderTransaction(request)
}

func (n *newOrderControllerImpl) Close() error {
	panic("implement me")
}

func makeNewOrderRequest(cmd []string, r *bufio.Reader) *model.NewOrderRequest {
	cIdString, wIdString, dIdString, mString := cmd[1], cmd[2], cmd[3], cmd[4]

	cId, _ := strconv.Atoi(cIdString)
	wId, _ := strconv.Atoi(wIdString)
	dId, _ := strconv.Atoi(dIdString)
	m, _ := strconv.Atoi(mString)

	return &model.NewOrderRequest{
		WId:              wId,
		DId:              dId,
		CId:              cId,
		NewOrderLineList: makeNewOrderLineList(m, r),
	}
}

func makeNewOrderLineList(m int, r *bufio.Reader) []*model.NewOrderLine {
	newOrderLineList := make([]*model.NewOrderLine, m)

	for i := 0; i < m; i++ {
		text, _ := r.ReadString('\n')
		orderLineSplit := strings.Split(strings.Trim(text, "\n"), ",")

		olIIdString, olSupplyWIdString, olQuantityString := orderLineSplit[0], orderLineSplit[1], orderLineSplit[2]
		olIId, _ := strconv.Atoi(olIIdString)
		olSupplyWId, _ := strconv.Atoi(olSupplyWIdString)
		olQuantity, _ := strconv.Atoi(olQuantityString)

		newOrderLine := &model.NewOrderLine{
			OlIId:       olIId,
			OlSupplyWId: olSupplyWId,
			OlQuantity:  olQuantity,
		}

		newOrderLineList = append(newOrderLineList, newOrderLine)
	}

	return newOrderLineList
}
