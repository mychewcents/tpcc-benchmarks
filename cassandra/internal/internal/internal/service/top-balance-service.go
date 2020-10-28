package service

import (
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/dao"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/datamodel/view"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/model"
	"io"
)

type TopBalanceService interface {
	ProcessTopBalanceTransaction() (*model.TopBalanceResponse, error)
	io.Closer
}

type topBalanceServiceImpl struct {
	c dao.CustomerDao
}

func NewTopBalanceService(cassandraSession *common.CassandraSession) TopBalanceService {
	return &topBalanceServiceImpl{
		c: dao.NewCustomerDao(cassandraSession),
	}
}

func (t *topBalanceServiceImpl) ProcessTopBalanceTransaction() (*model.TopBalanceResponse, error) {
	var topBalanceByWarehouse [10][10]*view.CustomerByBalanceView

	ch := make(chan [10]*view.CustomerByBalanceView, 10)
	for i := 1; i <= 10; i++ {
		go t.c.GetCustomerByTopNBalance(i, 10, ch)
	}

	for i := 0; i < 10; i++ {
		topBalanceByWarehouse[i] = <-ch
	}

	tbs := getTopNBalance(0, 9, topBalanceByWarehouse)
	var ci [10]*model.TopBalanceCustomerInfo

	for i, tb := range tbs {
		ci[i] = &model.TopBalanceCustomerInfo{
			CName: model.Name{
				FirstName:  tb.CName.FirstName,
				MiddleName: tb.CName.MiddleName,
				LastName:   tb.CName.LastName,
			},
			CBalance: tb.CBalance,
			WName:    tb.CWName,
			DName:    tb.CDName,
		}
	}

	return &model.TopBalanceResponse{CustomerInfoList: ci}, nil
}

func getTopNBalance(start int, end int, topBalanceByWarehouse [10][10]*view.CustomerByBalanceView) [10]*view.CustomerByBalanceView {
	if start == end {
		return topBalanceByWarehouse[start]
	}
	mid := (start + end) / 2
	tb1 := getTopNBalance(start, mid, topBalanceByWarehouse)
	tb2 := getTopNBalance(mid+1, end, topBalanceByWarehouse)

	return mergeTopNBalance(tb1, tb2)
}

func mergeTopNBalance(tb1 [10]*view.CustomerByBalanceView, tb2 [10]*view.CustomerByBalanceView) [10]*view.CustomerByBalanceView {
	var topNBalance [10]*view.CustomerByBalanceView

	j, k := 0, 0

	for i := 0; i < 10; i++ {
		if tb1[j].CBalance >= tb2[k].CBalance {
			topNBalance[i] = tb1[j]
			j++
		} else {
			topNBalance[i] = tb2[k]
			k++
		}
	}

	return topNBalance
}

func (t *topBalanceServiceImpl) Close() error {
	panic("implement me")
}
