package table

const (
	SWId       = "s_w_id"
	SIId       = "s_i_id"
	SQuantity  = "s_quantity"
	SIName     = "s_i_name"
	SIPrice    = "s_i_price"
	SIImId     = "s_i_im_id"
	SIData     = "s_i_data"
	SYtd       = "s_ytd"
	SOrderCnt  = "s_order_cnt"
	SRemoteCnt = "s_remote_cnt"
	SDist01    = "s_dist_01"
	SDist02    = "s_dist_02"
	SDist03    = "s_dist_03"
	SDist04    = "s_dist_04"
	SDist05    = "s_dist_05"
	SDist06    = "s_dist_06"
	SDist07    = "s_dist_07"
	SDist08    = "s_dist_08"
	SDist09    = "s_dist_09"
	SDist10    = "s_dist_10"
	SDate      = "s_date"
)

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
	SYtd       int64   `mapstructure:"s_ytd"`
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
	SDate      string  `mapstructure:"s_date"`
}
