package dao

import (
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
)

type WarehouseDao interface {
}

type warehouseDaoImpl struct {
	cassandraSession *common.CassandraSession
}

func NewWarehouseDao(cassandraSession *common.CassandraSession) WarehouseDao {
	return warehouseDaoImpl{cassandraSession: cassandraSession}
}
