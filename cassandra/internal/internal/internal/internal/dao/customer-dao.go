package dao

import (
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/datamodel/table"
	"log"
)

type CustomerDao interface {
	GetCustomerByKey(cWId int, cDId int, cId int, ch chan *table.CustomerTab)
	GetCustomerByTopNBalance(cWId int, n int, ch chan []*table.CustomerTab)
}

type customerDaoImpl struct {
	cluster *gocql.ClusterConfig
}

func NewCustomerDao(cluster *gocql.ClusterConfig) CustomerDao {
	return &customerDaoImpl{cluster: cluster}
}

func (c *customerDaoImpl) GetCustomerByKey(cWId int, cDId int, cId int, ch chan *table.CustomerTab) {
	session, err := c.cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	query := session.Query("SELECT * "+
		"from customer_tab "+
		"where c_w_id=? AND c_d_id=? and c_id=?", cWId, cDId, cId)

	result := make(map[string]interface{})
	if err := query.MapScan(result); err != nil {
		log.Fatalf("ERROR GetCustomerByKey error in query execution. cWId=%v, cDId=%v, cId=%v, err=%v\n", cWId, cDId, cId, err)
		return
	}

	ct, err := table.MakeCustomerTab(result)
	if err != nil {
		log.Fatalf("ERROR GetCustomerByKey error making customer. cWId=%v, cDId=%v, cId=%v, err=%v\n", cWId, cDId, cId, err)
		return
	}

	ch <- ct
}

func (c *customerDaoImpl) GetCustomerByTopNBalance(cWId int, n int, ch chan []*table.CustomerTab) {
	session, err := c.cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	query := session.Query("SELECT * "+
		"from customer_by_balance "+
		"where c_w_id=? limit ?", cWId, n)

	cts := make([]*table.CustomerTab, n)

	iter := query.Iter()
	defer iter.Close()

	for i, result := 0, make(map[string]interface{}); iter.MapScan(result); result, i = make(map[string]interface{}), i+1 {
		ct, err := table.MakeCustomerTab(result)
		if err != nil {
			log.Fatalf("ERROR GetCustomerByKey error making customer. cWId=%v, n=%v, err=%v\n", cWId, n, err)
			return
		}
		cts[i] = ct
	}

	ch <- cts
}
