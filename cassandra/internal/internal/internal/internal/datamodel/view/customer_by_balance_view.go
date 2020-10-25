package view

import (
	"github.com/mitchellh/mapstructure"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/datamodel/udt"
)

/*
CustomerByBalanceView maps to the customer_by_balance_view in cassandra
PrimaryKey((c_w_id), c_balance DESC, c_d_id ASC, c_id ASC)
*/
type CustomerByBalanceView struct {
	CWId     int      `mapstructure:"c_w_id"`
	CBalance float64  `mapstructure:"c_balance"`
	CDId     int      `mapstructure:"c_d_id"`
	CId      int      `mapstructure:"c_id"`
	CWName   string   `mapstructure:"c_w_name"`
	CDName   string   `mapstructure:"c_d_name"`
	CName    udt.Name `mapstructure:"c_name"`
}

func MakeCustomerByBalanceView(columns map[string]interface{}) (*CustomerByBalanceView, error) {
	var cv CustomerByBalanceView

	if err := mapstructure.Decode(columns, &cv); err != nil {
		return nil, err
	}

	return &cv, nil
}
