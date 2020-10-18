package dao

import (
	"github.com/mychewcents/ddbms-project/cassandra/internal/common"
)

type DistrictDao interface {
}

type districtDaoImpl struct {
	cassandraSession *common.CassandraSession
}

func NewDistrictDao(cassandraSession *common.CassandraSession) DistrictDao {
	return &districtDaoImpl{cassandraSession: cassandraSession}
}
