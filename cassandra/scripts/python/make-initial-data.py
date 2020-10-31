#!/usr/bin/env python
# coding: utf-8

import pandas as pd
import numpy as np
from cassandra import util
import time

# Warehouse
cols=["w_id", "w_name", "w_street_1", "w_street_2", "w_city", "w_state", "w_zip", "w_tax", "w_ytd"]
warehouse = pd.read_csv("assets/data/raw/warehouse.csv", header=None, names=cols)

w_address = "{street_1:'" + warehouse["w_street_1"] + "', street_2:'" +  warehouse["w_street_2"] + "', city:'" + warehouse["w_city"] + "', state:'" + warehouse["w_state"] + "', zip:" + warehouse["w_zip"].apply(str) + "}"

warehouse_tab = pd.concat([warehouse, w_address], axis=1)
warehouse_tab.columns = ['w_id', 'w_name', 'w_street_1', 'w_street_2', 'w_city','w_state', 'w_zip', 'w_tax','w_ytd', 'w_address']
warehouse_tab = warehouse_tab[['w_id', 'w_name', 'w_address', 'w_tax', 'w_ytd']]

warehouse_tab.to_csv("assets/data/processed/warehouse.csv", index=False)

# District
cols=["w_id", "d_id", "d_name", "d_street_1", "d_street_2", "d_city", "d_state", "d_zip", "d_tax", "d_ytd", "d_next_o_id"]
district = pd.read_csv("assets/data/raw/district.csv", header=None, names=cols)

d_address = "{street_1:'" + district["d_street_1"] + "', street_2:'" +  district["d_street_2"] + "', city:'" + district["d_city"] + "', state:'" + district["d_state"] + "', zip:" + district["d_zip"].apply(str) + "}"

district_tab = pd.concat([district, d_address], axis=1)
district_tab.columns = ['d_w_id', 'd_id', 'd_name', 'd_street_1', 'd_street_2', 'd_city','d_state"', 'd_zip', 'd_tax','d_ytd', 'd_next_o_id', 'd_address']
district_tab = district_tab[['d_w_id', 'd_id', 'd_name', 'd_address', 'd_tax','d_ytd']]

district_tab.to_csv("assets/data/processed/district.csv", index=False)

w_d = warehouse.set_index(["w_id"]).join(district.set_index(["w_id"]), on=["w_id"])
w_d.reset_index(inplace=True)
w_d = w_d[['w_id', 'w_name', 'w_tax', 'd_id', 'd_name', 'd_tax']]

# Customer
cols=["w_id", "d_id", "c_id", "c_first", "c_middle", "c_last", "c_street_1", "c_street_2", "c_city", "c_state", "c_zip", "c_phone",
     "c_since", "c_credit", "c_credit_lim", "c_discount", "c_balance", "c_ytd_payment", "c_payment_cnt", "c_delivery_cnt", "c_data"]
customer = pd.read_csv("assets/data/raw/customer.csv", header=None, names=cols)

c_join = w_d.set_index(["w_id", "d_id"]).join(customer.set_index(["w_id", "d_id"]), on=["w_id","d_id"])
c_join.reset_index(inplace=True)

c_join["c_address"] = "{street_1:'" + c_join["c_street_1"] + "', street_2:'" +  c_join["c_street_2"] + "', city:'" + c_join["c_city"] + "', state:'" + c_join["c_state"] + "', zip:" + c_join["c_zip"].apply(str) + "}"
c_join["c_name"] = "{first_name:'" + c_join["c_first"] + "', middle_name:'" +  c_join["c_middle"] + "', last_name:'" + c_join["c_last"] + "'}"

customer_tab = c_join
customer_tab.columns = ['c_w_id', 'c_d_id', 'c_w_name', 'c_w_tax', 'c_d_name', 'c_d_tax', 'c_id', 'c_first',
       'c_middle', 'c_last', 'c_street_1', 'c_street_2', 'c_city', 'c_state',
       'c_zip', 'c_phone', 'c_since', 'c_credit', 'c_credit_lim', 'c_discount',
       'c_balance', 'c_ytd_payment', 'c_payment_cnt', 'c_delivery_cnt',
       'c_data', 'c_address', 'c_name']
customer_tab = customer_tab[['c_w_id', 'c_d_id', 'c_id', 'c_w_name', 'c_w_tax', 'c_d_name', 'c_d_tax', 'c_name', 'c_address', 'c_phone', 'c_since', 'c_credit', 'c_credit_lim', 'c_discount', 'c_balance', 'c_ytd_payment', 'c_payment_cnt', 'c_delivery_cnt', 'c_data']]
customer_tab.to_csv("assets/data/processed/customer.csv", index=False)

# Stock
cols=["s_i_id", "s_i_name", "s_i_price", "s_i_im_id", "s_i_data"]
item = pd.read_csv("assets/data/raw/item.csv", header=None, names=cols)

cols=["s_w_id", "s_i_id", "s_quantity", "s_ytd", "s_order_cnt", "s_remote_cnt", "s_dist_01", "s_dist_02", "s_dist_03", "s_dist_04", "s_dist_05", "s_dist_06", "s_dist_07", "s_dist_08", "s_dist_09", "s_dist_10", "s_data"]
stock = pd.read_csv("assets/data/raw/stock.csv", header=None, names=cols)

i_s = item.set_index(["s_i_id"]).join(stock.set_index(["s_i_id"]), on=["s_i_id"])
i_s.reset_index(inplace=True)
i_s["s_ytd"] = i_s["s_ytd"].apply(int)

stock_tab = i_s[["s_w_id", "s_i_id", "s_quantity", "s_i_name", "s_i_price", "s_i_im_id", "s_i_data", "s_ytd", "s_order_cnt", "s_remote_cnt", "s_dist_01", "s_dist_02", "s_dist_03", "s_dist_04", "s_dist_05", "s_dist_06", "s_dist_07", "s_dist_08", "s_dist_09", "s_dist_10", "s_data"]]
stock_tab.to_csv("assets/data/processed/stock.csv", index=False)

# Order
cols=["w_id", "d_id", "o_id", "c_id", "o_carrier_id", "o_ol_cnt", "o_all_local", "o_entry_d"]
order = pd.read_csv("assets/data/raw/order.csv", header=None, names=cols)

order.sort_values(["c_id", "o_id"], inplace=True)

order['o_carrier_id'] = order['o_carrier_id'].fillna(-1).astype(np.int64)
order['o_carrier_id'].unique()

order["new_order_id"] = time.time()
order["new_order_id"] = order["new_order_id"] + order["o_id"]

order["new_order_id"] = order["new_order_id"].apply(util.uuid_from_time)

c_cols = c_join[["c_w_id", "c_d_id", "c_id", "c_name"]]
c_cols.columns = ["w_id", "d_id", "c_id", "c_name"]

o_join = order.set_index(["w_id", "d_id", "c_id"]).join(c_cols[["w_id", "d_id", "c_id", "c_name"]].set_index(["w_id", "d_id", "c_id"]), on=["w_id", "d_id", "c_id"])
o_join.reset_index(inplace=True)
o_join.sort_values(["w_id", "d_id", "c_id", "c_name"])

cols=["w_id", "d_id", "o_id", "ol_number", "i_id", "ol_delivery_d", "ol_amount", "ol_supply_w_id", "ol_quantity", "ol_dist_info"]
order_line = pd.read_csv("assets/data/raw/order-line.csv", header=None, names=cols)

order_line.sort_values(["w_id", "d_id", "o_id", "ol_number"])
order_line["ol_total_amount"] = order_line["ol_amount"] * order_line["ol_quantity"]

total_amount = order_line.groupby(["w_id", "d_id", "o_id"])["ol_total_amount"].sum()
total_amount_tab = pd.DataFrame(total_amount)
total_amount_tab.reset_index(inplace=True)
total_amount_tab.columns = ["w_id", "d_id", "o_id", "o_ol_total_amount"]

ol_delivery_d = order_line.groupby(["w_id", "d_id", "o_id"])["ol_delivery_d"].min()
ol_delivery_d_tab = pd.DataFrame(ol_delivery_d)
ol_delivery_d_tab.reset_index(inplace=True)
ol_delivery_d_tab.columns = ["w_id", "d_id", "o_id", "ol_delivery_d"]

o_ol_join = o_join.join(total_amount_tab.set_index(["w_id", "d_id", "o_id"]), on=["w_id", "d_id", "o_id"])
o_ol_join = o_ol_join.join(ol_delivery_d_tab.set_index(["w_id", "d_id", "o_id"]), on=["w_id", "d_id", "o_id"])

o_ol_join.columns = ['o_w_id', 'o_d_id', 'o_c_id', 'old_o_id', 'o_carrier_id', 'o_ol_count', 'o_all_local', 
                     'o_entry_d', 'o_id', 'o_c_name', 'o_ol_total_amount', 'ol_delivery_d']
o_ol_join['o_all_local'] = o_ol_join['o_all_local'].apply(lambda x: True if x == 1 else 0)

order_tab = o_ol_join[['o_w_id', 'o_d_id', 'o_id', 'o_c_id', 'o_c_name', 'o_carrier_id', 'ol_delivery_d', 'o_ol_count', 'o_ol_total_amount', 'o_all_local', 'o_entry_d']]
order_tab.to_csv("assets/data/processed/order.csv", index=False)

# Order-Line
o_ol_join_cols = o_ol_join[["o_w_id", "o_d_id", "old_o_id", "o_id"]]
o_ol_join_cols.columns = ["w_id", "d_id", "o_id", "new_o_id"]

ol_o_join = o_ol_join_cols.set_index(["w_id", "d_id", "o_id"]).join(order_line.set_index(["w_id", "d_id", "o_id"]), on=["w_id", "d_id", "o_id"])
ol_o_join.reset_index(inplace=True)

item_cols = item[["s_i_id", "s_i_name"]]
item_cols.columns = ["i_id", "i_name"]

item_ol_join = item_cols.set_index(["i_id"]).join(ol_o_join.set_index(["i_id"]), on=["i_id"])
item_ol_join.reset_index(inplace=True)

order_line_group_sum = item_ol_join[['w_id', 'd_id', 'o_id', 'i_id', 'ol_supply_w_id', 'ol_quantity']].groupby(by=['w_id', 'd_id', 'o_id', 'i_id', 'ol_supply_w_id'])
order_line_group_sum = order_line_group_sum.sum()

order_line_group_min = item_ol_join[['w_id', 'd_id', 'o_id', 'i_id', 'ol_supply_w_id', 'ol_dist_info']].groupby(by=['w_id', 'd_id', 'o_id', 'i_id', 'ol_supply_w_id'])
order_line_group_min = order_line_group_min.min()

order_line = order_line_group_sum.join(order_line_group_min)
order_line = order_line.reset_index()

order_line['ol_w_to_quantity'] = order_line['ol_supply_w_id'].apply(str) + ':' + order_line['ol_quantity'].apply(str)
order_line['ol_w_to_dist_info'] = order_line['ol_supply_w_id'].apply(str) + ':' + order_line['ol_dist_info'].apply(str)

order_line_group_list = order_line[['w_id', 'd_id', 'o_id', 'i_id', 'ol_w_to_quantity', 'ol_w_to_dist_info']].groupby(by=['w_id', 'd_id', 'o_id', 'i_id'])
order_line_list = order_line_group_list.agg(lambda x: list(x))

order_line_group_total_quantity = order_line[['w_id', 'd_id', 'o_id', 'i_id', 'ol_quantity']].groupby(by=['w_id', 'd_id', 'o_id', 'i_id'])
order_line_total = order_line_group_total_quantity.sum()

order_line_group_first = item_ol_join.groupby(by=['w_id', 'd_id', 'o_id', 'i_id'])
order_line_first = order_line_group_first.first()

ol = order_line_list.join(order_line_total)
ol = ol.join(order_line_first[['i_name', 'new_o_id', 'ol_number', 'ol_delivery_d', 'ol_amount', 'ol_supply_w_id', 'ol_dist_info', 'ol_total_amount']])
ol = ol.reset_index()

ol['ol_w_to_quantity'] = ol['ol_w_to_quantity'].apply(lambda x: "{" + ",".join(x) +"}")
ol['ol_w_to_dist_info'] = ol['ol_w_to_dist_info'].apply(lambda x: "{" + ",".join(x) +"}")

ol.columns = ['ol_w_id', 'ol_d_id', 'ol_old_o_id', 'ol_i_id', 'ol_w_to_quantity', 'ol_w_to_dist_info',
              'ol_quantity', 'ol_i_name', 'ol_o_id', 'ol_number', 'ol_delivery_d',
              'ol_amount', 'ol_supply_w_id', 'ol_dist_info', 'ol_total_amount']
order_line_tab = ol[['ol_w_id', 'ol_d_id', 'ol_o_id', 'ol_quantity', 'ol_number', 'ol_i_id', 'ol_i_name', 'ol_amount', 'ol_w_to_quantity', 'ol_w_to_dist_info']]
order_line_tab.to_csv("assets/data/processed/order-line.csv", index=False)

# Customer Item Pair
cols=["w_id", "d_id", "o_id", "c_id", "o_carrier_id", "o_ol_cnt", "o_all_local", "o_entry_d"]
order = pd.read_csv("assets/data/raw/order.csv", header=None, names=cols)
order.sort_values(["c_id", "o_id"], inplace=True)

cols=["w_id", "d_id", "o_id", "ol_number", "i_id", "ol_delivery_d", "ol_amount", "ol_supply_w_id", "ol_quantity", "ol_dist_info"]
order_line = pd.read_csv("assets/data/raw/order-line.csv", header=None, names=cols)
order_line.sort_values(["w_id", "d_id", "o_id", "ol_number"], inplace=True)

o_ol_join = order.set_index(["w_id", "d_id", "o_id"]).join(order_line.set_index(["w_id", "d_id", "o_id"]))
o_ol_join = o_ol_join.reset_index()
o_ol_join = o_ol_join[["w_id", "d_id", "o_id", "c_id", "i_id"]]

item_pair = o_ol_join.set_index(["w_id", "d_id", "o_id"]).join(o_ol_join.set_index(["w_id", "d_id", "o_id"]), lsuffix="l")

item_pair_distinct = item_pair[item_pair['i_idl'] < item_pair['i_id']]
item_pair_distinct = item_pair_distinct.reset_index()

item_pair_distinct['i_id_pair'] = "(" + item_pair_distinct['i_idl'].apply(str) + "," + item_pair_distinct['i_id'].apply(str) + ")"

item_pair_distinct_tab = item_pair_distinct[["w_id", "d_id", "c_id", "i_id_pair"]]
item_pair_distinct_tab.columns = ["c_w_id", "c_d_id", "c_id", "i_id_pair"]

import csv
item_pair_distinct_tab.to_csv("assets/data/processed/customer-item-order-pair.csv", index=False, quoting=csv.QUOTE_NONE, escapechar="\\")

