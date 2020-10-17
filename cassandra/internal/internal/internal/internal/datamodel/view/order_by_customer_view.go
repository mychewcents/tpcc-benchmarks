package view

import "time"

/*
OrderByCustomerView maps to the order_by_customer_view in cassandra
PrimaryKey((o_w_id, o_d_id), o_c_id ASC, o_id DESC)
*/
type OrderByCustomerView struct {
	OWId        int       `mapstructure:"o_w_id"`
	ODId        int       `mapstructure:"o_d_id"`
	OCId        int       `mapstructure:"o_c_id"`
	OId         int       `mapstructure:"o_id"`
	OEntryD     time.Time `mapstructure:"o_entry_d"`
	OCarrierId  int       `mapstructure:"o_carrier_id"`
	OlDeliveryD time.Time `mapstructure:"ol_delivery_d"`
}
