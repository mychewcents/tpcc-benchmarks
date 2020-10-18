package dao

import (
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/datamodel/table"
	"log"
)

type OrderLineDao interface {
	BatchInsertOrderLine(oltList []*table.OrderLineTab)
}

type orderLineDaoImpl struct {
	cluster *gocql.ClusterConfig
}

func NewOrderLineDao(cluster *gocql.ClusterConfig) OrderLineDao {
	return &orderLineDaoImpl{cluster: cluster}
}

func (o *orderLineDaoImpl) BatchInsertOrderLine(oltList []*table.OrderLineTab) {
	session, err := o.cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	batch := session.NewBatch(gocql.LoggedBatch)
	stmt := "INSERT INTO order_line_tab (ol_w_id, ol_d_id, ol_o_id, ol_quantity, ol_number, ol_i_id, ol_i_name, ol_amount, ol_supply_w_id, ol_dist_info) VALUES (?,?,?,?,?,?,?,?,?,?)"

	for _, ol := range oltList {
		batch.Query(stmt, ol.OlWId, ol.OlDId, ol.OlOId, ol.OlQuantity, ol.OlNumber, ol.OlIId, ol.OlIName, ol.OlAmount, ol.OlSupplyWId, ol.OlDistInfo)
	}

	err = session.ExecuteBatch(batch)
	if err != nil {
		log.Fatalf("ERROR BatchInsertNewOrderLine Error Executing batch err=%v", err)
	}
}
