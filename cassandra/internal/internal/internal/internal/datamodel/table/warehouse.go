package table

import (
	"github.com/mitchellh/mapstructure"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/datamodel/udt"
)

/*
WarehouseTab maps to the warehouse_tab in cassandra
PrimaryKey(w_id)
*/
type WarehouseTab struct {
	WId      int         `mapstructure:"w_id"`
	WName    string      `mapstructure:"w_name"`
	WAddress udt.Address `mapstructure:"w_address"`
	WTax     float32     `mapstructure:"w_tax"`
	WYtd     float64     `mapstructure:"w_ytd"`
}

func MakeWarehouseTab(columns map[string]interface{}) (*WarehouseTab, error) {
	var wt WarehouseTab

	if err := mapstructure.Decode(columns, &wt); err != nil {
		return nil, err
	}

	return &wt, nil
}
