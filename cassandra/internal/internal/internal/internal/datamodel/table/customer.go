package table

import (
	"github.com/mitchellh/mapstructure"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/datamodel/udt"
	"time"
)

/*
CustomerTab maps to the customer_tab in cassandra
PrimaryKey(c_w_id), c_d_id ASC, c_id ASC))
*/
type CustomerTab struct {
	CWId         int         `mapstructure:"c_w_id"`
	CWName       string      `mapstructure:"c_w_name"`
	CWTax        float32     `mapstructure:"c_w_tax"`
	CDId         int         `mapstructure:"c_d_id"`
	CDName       string      `mapstructure:"c_d_name"`
	CDTax        float32     `mapstructure:"c_d_tax"`
	CId          int         `mapstructure:"c_id"`
	CName        udt.Name    `mapstructure:"c_name"`
	CAddress     udt.Address `mapstructure:"c_address"`
	CPhone       string      `mapstructure:"c_phone"`
	CSince       time.Time   `mapstructure:"c_since"`
	CCredit      string      `mapstructure:"c_credit"`
	CCreditLim   float64     `mapstructure:"c_credit_lim"`
	CDiscount    float32     `mapstructure:"c_discount"`
	CBalance     float64     `mapstructure:"c_balance"`
	CYtdPayment  float64     `mapstructure:"c_ytd_payment"`
	CPaymentCnt  int         `mapstructure:"c_payment_cnt"`
	CDeliveryCnt int         `mapstructure:"c_delivery_cnt"`
	CData        string      `mapstructure:"c_data"`
}

func MakeCustomerTab(columns map[string]interface{}) (*CustomerTab, error) {
	var ct CustomerTab

	if err := mapstructure.Decode(columns, &ct); err != nil {
		return nil, err
	}

	return &ct, nil
}
