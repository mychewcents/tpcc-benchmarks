package dao

import (
	"github.com/gocql/gocql"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/datamodel/table"
	"log"
)

type StockDao interface {
	GetStockByKey(sWId int, sIId int, ch chan *table.StockTab)
	UpdateStockDaoCAS(stOld *table.StockTab, quantity int, isRemote bool, ch chan bool)
}

type stockDaoImpl struct {
	cluster *gocql.ClusterConfig
}

func NewStockDao(cluster *gocql.ClusterConfig) StockDao {
	return &stockDaoImpl{cluster: cluster}
}

func (s *stockDaoImpl) GetStockByKey(sWId int, sIId int, ch chan *table.StockTab) {
	session, err := s.cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	query := session.Query("SELECT * "+
		"from stock_tab "+
		"where s_w_id=? AND s_i_id=?", sWId, sIId)

	result := make(map[string]interface{})
	if err := query.MapScan(result); err != nil {
		log.Fatalf("ERROR GetStockByKey error in query execution. sWId=%v, sIId=%v, err=%v\n", sWId, sIId, err)
		return
	}

	st, err := table.MakeStockTab(result)
	if err != nil {
		log.Fatalf("ERROR GetStockByKey error making stock. sWId=%v, sIId=%v, err=%v\n", sWId, sIId, err)
		return
	}

	ch <- st
}

func (s *stockDaoImpl) UpdateStockDaoCAS(stOld *table.StockTab, quantity int, isRemote bool, ch chan bool) {
	session, err := s.cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	sQuantity := stOld.SQuantity - quantity
	if sQuantity < 10 {
		sQuantity = sQuantity + 100
	}
	sYtd := stOld.SYtd + quantity
	sOrderCnt := stOld.SOrderCnt + 1

	sRemoteCnt := stOld.SRemoteCnt
	if isRemote {
		sRemoteCnt++
	}

	query := session.Query("UPDATE stock_tab "+
		"SET s_quantity=?, s_ytd=?, s_order_cnt=?, s_remote_cnt=? "+
		"WHERE s_w_id=? and s_i_id=? "+
		"IF s_quantity=? AND s_ytd=? AND s_order_cnt=? AND s_remote_cnt=?", sQuantity, sYtd, sOrderCnt, sRemoteCnt,
		stOld.SWId, stOld.SIId,
		stOld.SQuantity, stOld.SYtd, stOld.SOrderCnt, stOld.SRemoteCnt)

	applied, err := query.ScanCAS(&sQuantity, &sYtd, &sOrderCnt, &sRemoteCnt)
	if err != nil {
		log.Fatalf("ERROR UpdateStockDaoCAS quering. err=%v\n", err)
		return
	}

	if !applied {
		log.Println("CAS Failure UpdateStockDaoCAS")
		stOld.SQuantity = sQuantity
		stOld.SYtd = sYtd
		stOld.SOrderCnt = sOrderCnt
		stOld.SRemoteCnt = sRemoteCnt

		s.UpdateStockDaoCAS(stOld, quantity, isRemote, ch)
	} else {
		ch <- true
	}
}
