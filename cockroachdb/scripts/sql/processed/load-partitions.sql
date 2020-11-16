IMPORT INTO ORDERS_WID_DID CSV DATA ('nodelocal://1/assets/processed/order/WID_DID.csv') WITH skip = '1', nullif='NULL';
IMPORT INTO ORDER_LINE_WID_DID CSV DATA ('nodelocal://1/assets/processed/orderline/WID_DID.csv') WITH skip = '1', nullif='NULL';
IMPORT INTO ORDER_ITEMS_CUSTOMERS_WID_DID CSV DATA ('nodelocal://1/assets/processed/itempairs/WID_DID.csv') WITH skip = '1', nullif='NULL';