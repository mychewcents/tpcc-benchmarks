package dao

import "github.com/gocql/gocql"

type StockDao interface {
}

type stockDaoImpl struct {
	cluster *gocql.ClusterConfig
}

func NewStockDao(cluster *gocql.ClusterConfig) StockDao {
	return stockDaoImpl{cluster: cluster}
}
