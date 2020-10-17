package dao

import (
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/datamodel/table"
	"strings"
)

type CustomerDao interface {
	GetCustomerByKey(cWId int, cDId int, cId int, columns []string) (*table.CustomerTab, error)
	GetCustomerByTopNBalance(cWId int, n int, columns []string) ([]*table.CustomerTab, error)
}

type customerDaoImpl struct {
	cluster *gocql.ClusterConfig
}

func NewCustomerDao(cluster *gocql.ClusterConfig) CustomerDao {
	return &customerDaoImpl{cluster: cluster}
}

func (c *customerDaoImpl) GetCustomerByKey(cWId int, cDId int, cId int, columns []string) (*table.CustomerTab, error) {
	session, _ := c.cluster.CreateSession()
	defer session.Close()

	query := session.Query("SELECT "+strings.Join(columns, ",")+
		" from customer_tab where c_w_id=? AND c_d_id=? and c_id=?", cWId, cDId, cId)

	result := make(map[string]interface{}, len(columns))
	if err := query.MapScan(result); err != nil {
		return nil, err
	}

	return table.MakeCustomerTab(result)
}

func (c *customerDaoImpl) GetCustomerByTopNBalance(cWId int, n int, columns []string) ([]*table.CustomerTab, error) {
	session, _ := c.cluster.CreateSession()
	defer session.Close()

	query := session.Query("SELECT "+strings.Join(columns, ",")+
		" from customer_by_balance where c_w_id=? limit ?", cWId, n)

	cts := make([]*table.CustomerTab, 0)

	iter := query.Iter()
	defer iter.Close()

	result := make(map[string]interface{}, len(columns))

	for ; iter.MapScan(result); result = make(map[string]interface{}, len(columns)) {
		ct, err := table.MakeCustomerTab(result)
		if err != nil {
			return nil, err
		}
		cts = append(cts, ct)
	}

	return cts, nil
}
