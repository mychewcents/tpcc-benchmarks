package dao

import "github.com/gocql/gocql"

type WarehouseDao interface {
}

type warehouseDaoImpl struct {
	cluster *gocql.ClusterConfig
}

func NewWarehouseDao(cluster *gocql.ClusterConfig) WarehouseDao {
	return warehouseDaoImpl{cluster: cluster}
}
