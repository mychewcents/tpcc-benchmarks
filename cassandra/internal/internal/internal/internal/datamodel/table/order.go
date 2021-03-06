package table

import (
	"github.com/gocql/gocql"
	"github.com/mitchellh/mapstructure"
	"github.com/mychewcents/tpcc-benchmarks/cassandra/internal/internal/internal/internal/datamodel/udt"
	"time"
)

/*
OrderTab maps to the order_tab in cassandra
PrimaryKey((o_w_id, o_d_id), o_id DESC)
*/
type OrderTab struct {
	OWId           int        `mapstructure:"o_w_id"`
	ODId           int        `mapstructure:"o_d_id"`
	OId            gocql.UUID `mapstructure:"o_id"`
	OCId           int        `mapstructure:"o_c_id"`
	OCName         udt.Name   `mapstructure:"o_c_name"`
	OCarrierId     int        `mapstructure:"o_carrier_id"`
	OlDeliveryD    time.Time  `mapstructure:"ol_delivery_d"`
	OOlCount       int        `mapstructure:"o_ol_count"`
	OOlTotalAmount float64    `mapstructure:"o_ol_total_amount"`
	OAllLocal      bool       `mapstructure:"o_all_local"`
	OEntryD        time.Time  `mapstructure:"o_entry_d"`
}

func MakeOrderTab(columns map[string]interface{}) (*OrderTab, error) {
	var ot OrderTab

	if err := mapstructure.Decode(columns, &ot); err != nil {
		return nil, err
	}

	return &ot, nil
}
