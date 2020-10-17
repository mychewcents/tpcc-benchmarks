package dao

import "github.com/gocql/gocql"

type OrderLineDao interface {
}

type orderLineDaoImpl struct {
	cluster *gocql.ClusterConfig
}

func NewOrderLineDao(cluster *gocql.ClusterConfig) OrderLineDao {
	return &orderLineDaoImpl{cluster: cluster}
}
