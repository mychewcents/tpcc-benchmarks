package dao

import "github.com/gocql/gocql"

type DistrictDao interface {
}

type districtDaoImpl struct {
	cluster *gocql.ClusterConfig
}

func NewDistrictDao(cluster *gocql.ClusterConfig) DistrictDao {
	return &districtDaoImpl{cluster: cluster}
}
