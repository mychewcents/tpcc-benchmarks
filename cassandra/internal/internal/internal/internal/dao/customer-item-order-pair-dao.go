package dao

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/datamodel/table"
	"log"
	"strings"
)

type CustomerItemOrderPairDao interface {
	GetCustomerItemOrderPairByCustomer(cWId int, cDId int, cId int) []*table.CustomerItemOrderPair
	GetCustomerItemOrderPairByItemPairList(cWId int, cDId int, itemPairList []string, ch chan []*table.CustomerItemOrderPair)
	BatchInsertCustomerItemOrderPair(ctList []*table.CustomerItemOrderPair, chComplete chan bool)
}

type customerItemOrderPairDaoImpl struct {
	cassandraSession *common.CassandraSession
}

func NewCustomerItemOrderPairDao(cassandraSession *common.CassandraSession) CustomerItemOrderPairDao {
	return &customerItemOrderPairDaoImpl{cassandraSession: cassandraSession}
}

func (c *customerItemOrderPairDaoImpl) GetCustomerItemOrderPairByCustomer(cWId int, cDId int, cId int) []*table.CustomerItemOrderPair {
	query := c.cassandraSession.ReadSession.Query("SELECT * "+
		"from customer_item_order_pair_tab "+
		"where c_w_id=? AND c_d_id=? and c_id=?", cWId, cDId, cId)

	cts := make([]*table.CustomerItemOrderPair, 0)

	iter := query.Iter()
	defer iter.Close()

	for result := make(map[string]interface{}); iter.MapScan(result); result = make(map[string]interface{}) {
		ct, err := table.MakeCustomerItemOrderPair(result)
		if err != nil {
			log.Fatalf("ERROR GetCustomerItemOrderPairByCustomer error making orderLine. cWId=%v, cDId=%v, cId=%v, err=%v\n", cWId, cDId, cId, err)
			return nil
		}
		cts = append(cts, ct)
	}

	return cts
}

func (c *customerItemOrderPairDaoImpl) GetCustomerItemOrderPairByItemPairList(cWId int, cDId int, itemPairList []string, ch chan []*table.CustomerItemOrderPair) {
	stmt := fmt.Sprintf("select * from customer_item_order_pair_by_item_pair "+
		"where c_w_id=%d AND c_d_id=%d AND i_id_pair IN (%s)", cWId, cDId, strings.Join(itemPairList, ","))

	query := c.cassandraSession.ReadSession.Query(stmt)

	cts := make([]*table.CustomerItemOrderPair, 0)

	iter := query.Iter()
	defer iter.Close()

	for result := make(map[string]interface{}); iter.MapScan(result); result = make(map[string]interface{}) {
		ct, err := table.MakeCustomerItemOrderPair(result)
		if err != nil {
			log.Fatalf("ERROR GetCustomerItemOrderPairByItemPairList error making orderLine. cWId=%v, cDId=%v, itemPairList=%v, err=%v\n", cWId, cDId, strings.Join(itemPairList, ","), err)
			return
		}
		cts = append(cts, ct)
	}

	ch <- cts
}

func (c *customerItemOrderPairDaoImpl) BatchInsertCustomerItemOrderPair(ctList []*table.CustomerItemOrderPair, chComplete chan bool) {
	batch := c.cassandraSession.WriteSession.NewBatch(gocql.LoggedBatch)
	stmt := "INSERT INTO customer_item_order_pair_tab (c_w_id, c_d_id, c_id, i_id_pair) VALUES (?,?,?,(?,?))"

	for _, ct := range ctList {
		batch.Query(stmt, ct.CWId, ct.CDId, ct.CId, ct.IId1, ct.IId2)
	}

	err := c.cassandraSession.WriteSession.ExecuteBatch(batch)
	if err != nil {
		log.Fatalf("ERROR BatchInsertCustomerItemOrderPair Error Executing batch err=%v", err)
	}

	chComplete <- true
}
