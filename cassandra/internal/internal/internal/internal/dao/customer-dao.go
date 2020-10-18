package dao

import (
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/datamodel/table"
	"log"
)

type CustomerDao interface {
	GetCustomerByKey(cWId int, cDId int, cId int, ch chan *table.CustomerTab)
	GetCustomerByTopNBalance(cWId int, n int, ch chan []*table.CustomerTab)
	UpdateCustomerCAS(ctOld *table.CustomerTab, payment float64, ch chan bool)
}

type customerDaoImpl struct {
	cassandraSession *common.CassandraSession
}

func NewCustomerDao(cassandraSession *common.CassandraSession) CustomerDao {
	return &customerDaoImpl{cassandraSession: cassandraSession}
}

func (c *customerDaoImpl) GetCustomerByKey(cWId int, cDId int, cId int, ch chan *table.CustomerTab) {

	query := c.cassandraSession.ReadSession.Query("SELECT * "+
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
	query := c.cassandraSession.ReadSession.Query("SELECT * "+
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

func (c *customerDaoImpl) UpdateCustomerCAS(ctOld *table.CustomerTab, payment float64, ch chan bool) {

	cBalance := ctOld.CBalance - payment
	cYtdPayment := ctOld.CYtdPayment + payment
	cPaymentCnt := ctOld.CPaymentCnt - 1

	query := c.cassandraSession.WriteSession.Query("UPDATE customer_tab "+
		"SET c_balance=?, c_ytd_payment=? c_payment_cnt=? "+
		"WHERE c_w_id=? AND c_d_id=? AND c_id=? "+
		"IF c_balance=?, c_ytd_payment=? c_payment_cnt=?", cBalance, cYtdPayment, cPaymentCnt,
		ctOld.CWId, ctOld.CDId, ctOld.CId,
		ctOld.CBalance, ctOld.CYtdPayment, ctOld.CPaymentCnt)

	applied, err := query.ScanCAS(&cBalance, &cYtdPayment, &cPaymentCnt)
	if err != nil {
		log.Fatalf("ERROR UpdateWarehouseCAS quering. err=%v\n", err)
		return
	}

	if !applied {
		log.Println("CAS Failure UpdateWarehouseCAS")
		ctOld.CBalance = cBalance
		ctOld.CYtdPayment = cYtdPayment
		ctOld.CPaymentCnt = cPaymentCnt

		c.UpdateCustomerCAS(ctOld, payment, ch)
	} else {
		ch <- true
	}
}
