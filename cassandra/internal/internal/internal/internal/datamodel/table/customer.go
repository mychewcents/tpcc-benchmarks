package table

import (
	"github.com/mitchellh/mapstructure"
	"github.com/mychewcents/ddbms-project/cassandra/internal/internal/internal/internal/datamodel/udt"
	"time"
)

const (
	CWId         = "c_w_id"
	CWName       = "c_w_name"
	CWTax        = "c_w_tax"
	CDId         = "c_d_id"
	CDName       = "c_d_name"
	CDTax        = "c_d_tax"
	CId          = "c_id"
	CName        = "c_name"
	CAddress     = "c_address"
	CPhone       = "c_phone"
	CSince       = "c_since"
	CCredit      = "c_credit"
	CCreditLim   = "c_credit_lim"
	CDiscount    = "c_discount"
	CBalance     = "c_balance"
	CYTDPayment  = "c_ytd_payment"
	CPaymentCnt  = "c_payment_cnt"
	CDeliveryCnt = "c_delivery_cnt"
	CData        = "c_data"
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
