package dao

import (
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/common"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/internal/internal/datamodel/table"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/internal/internal/datamodel/view"
	"log"
)

type CustomerDao interface {
	GetCustomerByKey(cWId int, cDId int, cId int, ch chan *table.CustomerTab)
	GetCustomerByTopNBalance(cWId int, n int, ch chan [10]*view.CustomerByBalanceView)
	UpdateCustomerPaymentCAS(ctOld *table.CustomerTab, payment float64, ch chan bool)
	UpdateCustomerDeliveryCAS(ctOld *table.CustomerTab, olAmount float64)
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

func (c *customerDaoImpl) GetCustomerByTopNBalance(cWId int, n int, ch chan [10]*view.CustomerByBalanceView) {
	query := c.cassandraSession.ReadSession.Query("SELECT * "+
		"from customer_by_balance "+
		"where c_w_id=? limit ?", cWId, n)

	var cts [10]*view.CustomerByBalanceView

	iter := query.Iter()
	defer iter.Close()

	for i, result := 0, make(map[string]interface{}); iter.MapScan(result); result, i = make(map[string]interface{}), i+1 {
		ct, err := view.MakeCustomerByBalanceView(result)
		if err != nil {
			log.Fatalf("ERROR GetCustomerByTopNBalance error making customer. cWId=%v, n=%v, err=%v\n", cWId, n, err)
			return
		}
		cts[i] = ct
	}

	ch <- cts
}

func (c *customerDaoImpl) UpdateCustomerPaymentCAS(ctOld *table.CustomerTab, payment float64, ch chan bool) {

	cBalance := ctOld.CBalance - payment
	cYtdPayment := ctOld.CYtdPayment + payment
	cPaymentCnt := ctOld.CPaymentCnt + 1

	query := c.cassandraSession.WriteSession.Query("UPDATE customer_tab "+
		"SET c_balance=?, c_ytd_payment=?, c_payment_cnt=? "+
		"WHERE c_w_id=? AND c_d_id=? AND c_id=? "+
		"IF c_balance=? AND c_ytd_payment=? AND c_payment_cnt=?", cBalance, cYtdPayment, cPaymentCnt,
		ctOld.CWId, ctOld.CDId, ctOld.CId,
		ctOld.CBalance, ctOld.CYtdPayment, ctOld.CPaymentCnt)

	applied, err := query.ScanCAS(&cBalance, &cYtdPayment, &cPaymentCnt)
	if err != nil {
		log.Fatalf("ERROR UpdateCustomerPaymentCAS quering. err=%v\n", err)
		return
	}

	if !applied {
		log.Println("CAS Failure UpdateCustomerPaymentCAS")
		ctOld.CBalance = cBalance
		ctOld.CYtdPayment = cYtdPayment
		ctOld.CPaymentCnt = cPaymentCnt

		c.UpdateCustomerPaymentCAS(ctOld, payment, ch)
	} else {
		ch <- true
	}
}

func (c *customerDaoImpl) UpdateCustomerDeliveryCAS(ctOld *table.CustomerTab, olAmount float64) {
	cBalance := ctOld.CBalance + olAmount
	cDeliveryCnt := ctOld.CDeliveryCnt + 1

	query := c.cassandraSession.WriteSession.Query("UPDATE customer_tab "+
		"SET c_balance=?, c_delivery_cnt=? "+
		"WHERE c_w_id=? AND c_d_id=? AND c_id=? "+
		"IF c_balance=? AND c_delivery_cnt=?", cBalance, cDeliveryCnt,
		ctOld.CWId, ctOld.CDId, ctOld.CId,
		ctOld.CBalance, ctOld.CDeliveryCnt)

	applied, err := query.ScanCAS(&cBalance, &cDeliveryCnt)
	if err != nil {
		log.Fatalf("ERROR UpdateCustomerDeliveryCAS quering. err=%v\n", err)
		return
	}

	if !applied {
		log.Println("CAS Failure UpdateCustomerDeliveryCAS")
		ctOld.CBalance = cBalance
		ctOld.CDeliveryCnt = cDeliveryCnt

		c.UpdateCustomerDeliveryCAS(ctOld, olAmount)
	}
}
