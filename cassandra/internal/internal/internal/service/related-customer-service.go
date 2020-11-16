package service

import (
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/common"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/internal/internal/dao"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/internal/internal/datamodel/table"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/internal/model"
	"io"
)

type RelatedCustomerService interface {
	ProcessRelatedCustomerTransaction(request *model.RelatedCustomerRequest) (*model.RelatedCustomerResponse, error)
	io.Closer
}

type relatedCustomerServiceImpl struct {
	c dao.CustomerItemOrderPairDao
}

func NewRelatedCustomerService(cassandraSession *common.CassandraSession) RelatedCustomerService {
	return &relatedCustomerServiceImpl{
		c: dao.NewCustomerItemOrderPairDao(cassandraSession),
	}
}

func (r *relatedCustomerServiceImpl) ProcessRelatedCustomerTransaction(request *model.RelatedCustomerRequest) (*model.RelatedCustomerResponse, error) {
	customerIdentifier := request.CustomerIdentifier
	cts := r.c.GetCustomerItemOrderPairByCustomer(customerIdentifier.CWId, customerIdentifier.CDId, customerIdentifier.CId)

	cis := r.getRelatedCustomers(customerIdentifier.CWId, makeItemPairTupleList(cts))

	return &model.RelatedCustomerResponse{
		CustomerIdentifier:         customerIdentifier,
		RelatedCustomerIdentifiers: cis,
	}, nil
}

func (r *relatedCustomerServiceImpl) getRelatedCustomers(cWId int, itemPairTupleList []string) []*model.CustomerIdentifier {
	cis := make([]*model.CustomerIdentifier, 0)
	ch := make(chan []*table.CustomerItemOrderPair, 90)

	for wId := 1; wId <= 10; wId++ {
		if wId != cWId {
			for dId := 1; dId <= 10; dId++ {
				go r.c.GetCustomerItemOrderPairByItemPairList(wId, dId, itemPairTupleList, ch)
			}
		}
	}

	for i := 0; i < 90; i++ {
		cts := <-ch
		cis = append(cis, makeCustomerIdentifier(cts)...)
	}

	return cis
}

func makeCustomerIdentifier(cts []*table.CustomerItemOrderPair) []*model.CustomerIdentifier {
	cis := make([]*model.CustomerIdentifier, 0)

	// Required because customers may have multiple item-pairs in common
	ciMap := make(map[int]map[int]map[int]bool) // Map [cWId][cDId][cId]

	for _, ct := range cts {
		if !ciMap[ct.CWId][ct.CDId][ct.CId] {

			if ciMap[ct.CWId] == nil {
				ciMap[ct.CWId] = make(map[int]map[int]bool)
			}
			if ciMap[ct.CWId][ct.CDId] == nil {
				ciMap[ct.CWId][ct.CDId] = make(map[int]bool)
			}
			ciMap[ct.CWId][ct.CDId][ct.CId] = true

			ci := &model.CustomerIdentifier{
				CWId: ct.CWId,
				CDId: ct.CDId,
				CId:  ct.CId,
			}

			cis = append(cis, ci)
		}
	}

	return cis
}

func makeItemPairTupleList(cts []*table.CustomerItemOrderPair) []string {
	ips := make([]string, len(cts))
	for i, ct := range cts {
		ips[i] = ct.GetItemIdPair()
	}
	return ips
}

func (r *relatedCustomerServiceImpl) Close() error {
	panic("implement me")
}
