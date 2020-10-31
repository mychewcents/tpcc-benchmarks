package table

import (
	"github.com/mitchellh/mapstructure"
	"strconv"
)

type CustomerItemOrderPair struct {
	CWId int `mapstructure:"c_w_id"`
	CDId int `mapstructure:"c_d_id"`
	CId  int `mapstructure:"c_id"`
	IId1 int `mapstructure:"i_id_pair[0]"`
	IId2 int `mapstructure:"i_id_pair[1]"`
}

func (c *CustomerItemOrderPair) GetItemIdPair() string {
	return "(" + strconv.Itoa(c.IId1) + ", " + strconv.Itoa(c.IId2) + ")"
}

func MakeCustomerItemOrderPair(columns map[string]interface{}) (*CustomerItemOrderPair, error) {
	var ct CustomerItemOrderPair

	if err := mapstructure.Decode(columns, &ct); err != nil {
		return nil, err
	}

	return &ct, nil
}
