package dao

import "github.com/gocql/gocql"

type OrderDao interface {
}

type orderDaoImpl struct {
	cluster *gocql.ClusterConfig
}

func NewOrderDao(cluster *gocql.ClusterConfig) OrderDao {
	return &orderDaoImpl{cluster: cluster}
}
