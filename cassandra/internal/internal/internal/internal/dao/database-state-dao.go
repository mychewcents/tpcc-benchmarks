package dao

import (
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
	"log"
)

type DatabaseStateDao interface {
	GetWarehouseState() (sumWYTD float64)
	GetDistrictState() (sumDYTD float64)
	GetCustomerState() (sumCBalance float64, sumCYTDPayment float64, sumCPaymentCnt int64, sumCDeliveryCnt int64)
	GetOrderState() (maxOId gocql.UUID, sumOOlCnt int64)
	GetOrderLineState() (sumOlAmount float64, sumOlQuantity int64)
	GetStockState() (sumSQuantity int64, sumSYTD int64, sumSOrderCnt int64, sumSRemoteCnt int64)
}

type databaseStateDaoImpl struct {
	cassandraSession *common.CassandraSession
}

func NewDatabaseStateDao(cassandraSession *common.CassandraSession) DatabaseStateDao {
	return &databaseStateDaoImpl{cassandraSession: cassandraSession}
}

func (d *databaseStateDaoImpl) GetWarehouseState() (sumWYTD float64) {
	query := d.cassandraSession.ReadSession.Query("select sum(w_ytd) from warehouse_tab")
	err := query.Scan(sumWYTD)
	if err != nil {
		log.Printf("ERROR GetWarehouseState err=%v\n", err)
		return
	}
	return
}

func (d *databaseStateDaoImpl) GetDistrictState() (sumDYTD float64) {
	query := d.cassandraSession.ReadSession.Query("select sum(d_ytd) from district_tab")
	err := query.Scan(sumDYTD)
	if err != nil {
		log.Printf("ERROR GetDistrictState err=%v\n", err)
		return
	}
	return
}

func (d *databaseStateDaoImpl) GetCustomerState() (sumCBalance float64, sumCYTDPayment float64, sumCPaymentCnt int64, sumCDeliveryCnt int64) {
	query := d.cassandraSession.ReadSession.Query("select sum(c_balance), sum(c_ytd_payment), sum(c_payment_cnt), sum(c_delivery_cnt) from customer_tab")
	err := query.Scan(sumCBalance, sumCYTDPayment, sumCPaymentCnt, sumCDeliveryCnt)
	if err != nil {
		log.Printf("ERROR GetCustomerState err=%v\n", err)
		return
	}
	return
}

func (d *databaseStateDaoImpl) GetOrderState() (maxOId gocql.UUID, sumOOlCnt int64) {
	query := d.cassandraSession.ReadSession.Query("select max(o_id), sum(o_ol_count) from order_tab")
	err := query.Scan(maxOId, sumOOlCnt)
	if err != nil {
		log.Printf("ERROR GetCustomerState err=%v\n", err)
		return
	}
	return
}

func (d *databaseStateDaoImpl) GetOrderLineState() (sumOlAmount float64, sumOlQuantity int64) {
	query := d.cassandraSession.ReadSession.Query("select sum(ol_amount), sum(ol_quantity) from order_line_tab")
	err := query.Scan(sumOlAmount, sumOlQuantity)
	if err != nil {
		log.Printf("ERROR GetCustomerState err=%v\n", err)
		return
	}
	return
}

func (d *databaseStateDaoImpl) GetStockState() (sumSQuantity int64, sumSYTD int64, sumSOrderCnt int64, sumSRemoteCnt int64) {
	query := d.cassandraSession.ReadSession.Query("select sum(s_quantity), sum(s_ytd), sum(s_order_cnt), sum(s_remote_cnt) from stock_tab")
	err := query.Scan(sumSQuantity, sumSYTD, sumSOrderCnt, sumSRemoteCnt)
	if err != nil {
		log.Printf("ERROR GetCustomerState err=%v\n", err)
		return
	}
	return
}
