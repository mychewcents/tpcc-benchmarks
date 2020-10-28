package controller

import (
	"fmt"
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/handler"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/service"
	"strconv"
)

type PopularItemController interface {
	handler.TransactionHandler
}

type popularItemControllerImpl struct {
	s service.PopularItemService
}

func NewPopularItemController(cassandraSession *common.CassandraSession) PopularItemController {
	return &popularItemControllerImpl{
		s: service.NewPopularItemService(cassandraSession),
	}
}

func (p *popularItemControllerImpl) HandleTransaction(cmd []string) {
	request := makePopularItemRequest(cmd)
	response, _ := p.s.ProcessPopularItemService(request)
	printPopularItemResponse(response)
}

func makePopularItemRequest(cmd []string) *model.PopularItemRequest {
	wId, _ := strconv.Atoi(cmd[1])
	dId, _ := strconv.Atoi(cmd[2])
	l, _ := strconv.Atoi(cmd[3])

	return &model.PopularItemRequest{
		WId:            wId,
		DId:            dId,
		NoOfLastOrders: l,
	}
}

func printPopularItemResponse(r *model.PopularItemResponse) {
	//fmt.Println(r)
	fmt.Sprintf("1.district identifier, & Number of last orders to be examined L:%v\n", response)
    fmt.Sprintf("3.For each order number x in S:%v\n", response.OrderItemInfoList)
    fmt.Sprintf("4.The percentage of examined orders that contain each popular item:%v\n", response.PopularItemStatList)
}

func (p *popularItemControllerImpl) Close() error {
	panic("implement me")
}
