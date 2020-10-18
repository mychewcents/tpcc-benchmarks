package dao

import (
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/datamodel/table"
	"log"
)

type OrderDao interface {
	InsertOrder(ot *table.OrderTab)
}

type orderDaoImpl struct {
	cluster *gocql.ClusterConfig
}

func NewOrderDao(cluster *gocql.ClusterConfig) OrderDao {
	return &orderDaoImpl{cluster: cluster}
}

func (o *orderDaoImpl) InsertOrder(ot *table.OrderTab) {
	session, err := o.cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	stmt := "INSERT INTO " +
		"order_tab (o_w_id, o_d_id, o_id, o_c_id, o_c_name, o_carrier_id, ol_delivery_d, o_ol_count, o_ol_total_amount, o_all_local, o_entry_d) " +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?)"

	query := session.Query(stmt, ot.OWId, ot.ODId, ot.OId, ot.OCId, ot.OCName.GetNameString(), ot.OCarrierId, ot.OlDeliveryD,
		ot.OOlCount, ot.OOlTotalAmount, ot.OAllLocal, ot.OEntryD)

	err = query.Exec()
	if err != nil {
		log.Fatalf("InsertOrder. ot=%v, err%v", ot, err)
	}
}
