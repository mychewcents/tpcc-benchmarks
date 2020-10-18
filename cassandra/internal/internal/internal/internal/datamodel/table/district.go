package table

import "github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/datamodel/udt"

const (
	DWId     = "d_w_id"
	DId      = "d_id"
	DName    = "d_name"
	DAddress = "d_address"
	DTax     = "d_tax"
	DYtd     = "d_ytd"
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
