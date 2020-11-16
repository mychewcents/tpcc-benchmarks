package controller

import (
	"fmt"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/common"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/internal/service"
	"log"
	"os"
)

type DatabaseStateController interface {
	SaveDatabaseState(path string, experimentNo int)
}

type databaseStateControllerImpl struct {
	d service.DatabaseStateService
}

func NewDatabaseStateController(cassandraSession *common.CassandraSession) DatabaseStateController {
	return &databaseStateControllerImpl{d: service.NewDatabaseStateService(cassandraSession)}
}

func (d *databaseStateControllerImpl) SaveDatabaseState(path string, experimentNo int) {
	state, _ := d.d.GetDatabaseState()

	fileName := fmt.Sprintf("%v/experiment_%v.csv", path, experimentNo)
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("ERROR saving database state, err=%v", err)
	}
	defer file.Close()

	databaseState := fmt.Sprintf("%d,%.2f,%.2f,'%s',%.2f,%.2f,%d,%d,%v,%d,%.2f,%d,%d,%d,%d,%d", experimentNo, state.SumWYTD, state.SumDYTD, "N/A",
		state.SumCBalance, state.SumCYTDPayment, state.SumCPaymentCnt, state.SumCDeliveryCnt,
		state.MaxOId, state.SumOOlCnt, state.SumOlAmount, state.SumOlQuantity,
		state.SumSQuantity, state.SumSYTD, state.SumSOrderCnt, state.SumSRemoteCnt)

	file.WriteString(databaseState)
}
