use cassandra;
COPY stock_tab (s_w_id,s_i_id,s_quantity,s_i_name,s_i_price,s_i_im_id,s_i_data,s_ytd,s_order_cnt,s_remote_cnt,s_dist_01,s_dist_02,s_dist_03,s_dist_04,s_dist_05,s_dist_06,s_dist_07,s_dist_08,s_dist_09,s_dist_10,s_data) FROM 'assets/data/processed/stock.csv' WITH HEADER = true ;
COPY warehouse_tab (w_id, w_name, w_address, w_tax, w_ytd) FROM 'assets/data/processed/warehouse.csv' WITH HEADER = true ;
COPY district_tab (d_w_id, d_id, d_name, d_address, d_tax,d_ytd) FROM 'assets/data/processed/district.csv' WITH HEADER = true ;
COPY customer_tab (c_w_id, c_d_id, c_id, c_w_name, c_w_tax, c_d_name, c_d_tax, c_name, c_address, c_phone, c_since, c_credit, c_credit_lim, c_discount, c_balance, c_ytd_payment, c_payment_cnt, c_delivery_cnt, c_data) FROM 'assets/data/processed/customer.csv' WITH HEADER = true ;
COPY order_tab (o_w_id, o_d_id, o_id, o_c_id, o_c_name, o_carrier_id, ol_delivery_d, o_ol_count, o_ol_total_amount, o_all_local, o_entry_d) FROM 'assets/data/processed/order.csv' WITH HEADER = true ;
COPY order_line_tab (ol_w_id, ol_d_id, ol_o_id, ol_quantity, ol_number, ol_i_id, ol_i_name, ol_amount, ol_w_to_quantity, ol_w_to_dist_info) FROM 'assets/data/processed/order-line.csv' WITH HEADER = true ;
COPY customer_item_order_pair_tab (c_w_id, c_d_id , c_id , i_id_pair) FROM 'assets/data/processed/customer-item-order-pair.csv' WITH HEADER = true ;
