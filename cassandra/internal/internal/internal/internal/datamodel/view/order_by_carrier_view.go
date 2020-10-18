package view

/*
OrderByCarrierView maps to the order_by_carrier_view in cassandra
PrimaryKey((o_w_id, o_d_id), o_carrier_id ASC, o_id ASC)
*/
type OrderByCarrierView struct {
	OWId       int `mapstructure:"o_w_id"`
	ODId       int `mapstructure:"o_d_id"`
	OCarrierId int `mapstructure:"o_carrier_id"`
	OId        int `mapstructure:"o_id"`
}
