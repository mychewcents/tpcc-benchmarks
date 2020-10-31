package table

import "github.com/mitchellh/mapstructure"

/*
StockTab maps to the stock_tab in cassandra
PrimaryKey((s_w_id, s_i_id), s_quantity DESC)
*/
type StockTab struct {
	SWId       int     `mapstructure:"s_w_id"`
	SIId       int     `mapstructure:"s_i_id"`
	SQuantity  int     `mapstructure:"s_quantity"`
	SIName     string  `mapstructure:"s_i_name"`
	SIPrice    float32 `mapstructure:"s_i_price"`
	SIImId     int     `mapstructure:"s_i_im_id"`
	SIData     string  `mapstructure:"s_i_data"`
	SYtd       int     `mapstructure:"s_ytd"`
	SOrderCnt  int     `mapstructure:"s_order_cnt"`
	SRemoteCnt int     `mapstructure:"s_remote_cnt"`
	SDist01    string  `mapstructure:"s_dist_01"`
	SDist02    string  `mapstructure:"s_dist_02"`
	SDist03    string  `mapstructure:"s_dist_03"`
	SDist04    string  `mapstructure:"s_dist_04"`
	SDist05    string  `mapstructure:"s_dist_05"`
	SDist06    string  `mapstructure:"s_dist_06"`
	SDist07    string  `mapstructure:"s_dist_07"`
	SDist08    string  `mapstructure:"s_dist_08"`
	SDist09    string  `mapstructure:"s_dist_09"`
	SDist10    string  `mapstructure:"s_dist_10"`
	SData      string  `mapstructure:"s_data"`
}

func (st *StockTab) GetSDist(dId int) string {
	switch dId {
	case 1:
		return st.SDist01
	case 2:
		return st.SDist02
	case 3:
		return st.SDist03
	case 4:
		return st.SDist04
	case 5:
		return st.SDist05
	case 6:
		return st.SDist06
	case 7:
		return st.SDist07
	case 8:
		return st.SDist08
	case 9:
		return st.SDist09
	case 10:
		return st.SDist10
	}

	return ""
}

func MakeStockTab(columns map[string]interface{}) (*StockTab, error) {
	var st StockTab

	if err := mapstructure.Decode(columns, &st); err != nil {
		return nil, err
	}

	return &st, nil
}
