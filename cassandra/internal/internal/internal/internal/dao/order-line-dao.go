package dao

import (
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/datamodel/table"
	"log"
)

type OrderLineDao interface {
	BatchInsertOrderLine(oltList []*table.OrderLineTab)
}

type orderLineDaoImpl struct {
	cassandraSession *common.CassandraSession
}

func NewOrderLineDao(cassandraSession *common.CassandraSession) OrderLineDao {
	return &orderLineDaoImpl{cassandraSession: cassandraSession}
}

func (o *orderLineDaoImpl) BatchInsertOrderLine(oltList []*table.OrderLineTab) {
	batch := o.cassandraSession.WriteSession.NewBatch(gocql.LoggedBatch)
	stmt := "INSERT INTO order_line_tab (ol_w_id, ol_d_id, ol_o_id, ol_quantity, ol_number, ol_i_id, ol_i_name, ol_amount, ol_supply_w_id, ol_dist_info) VALUES (?,?,?,?,?,?,?,?,?,?)"

	for _, ol := range oltList {
		batch.Query(stmt, ol.OlWId, ol.OlDId, ol.OlOId, ol.OlQuantity, ol.OlNumber, ol.OlIId, ol.OlIName, ol.OlAmount, ol.OlSupplyWId, ol.OlDistInfo)
	}

	err := o.cassandraSession.WriteSession.ExecuteBatch(batch)
	if err != nil {
		log.Fatalf("ERROR BatchInsertNewOrderLine Error Executing batch err=%v", err)
	}
}
