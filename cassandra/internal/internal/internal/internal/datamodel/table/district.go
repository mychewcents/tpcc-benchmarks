package table

import (
	"github.com/mitchellh/mapstructure"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/internal/internal/datamodel/udt"
)

/*
DistrictTab maps to the district_tab in cassandra
PrimaryKey(d_w_id, d_id)
*/
type DistrictTab struct {
	DWId     int         `mapstructure:"d_w_id"`
	DId      int         `mapstructure:"d_id"`
	DName    string      `mapstructure:"d_name"`
	DAddress udt.Address `mapstructure:"d_address"`
	DTax     float32     `mapstructure:"d_tax"`
	DYtd     float64     `mapstructure:"d_ytd"`
}

func MakeDistrictTab(columns map[string]interface{}) (*DistrictTab, error) {
	var dt DistrictTab

	if err := mapstructure.Decode(columns, &dt); err != nil {
		return nil, err
	}

	return &dt, nil
}
