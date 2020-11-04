IMPORT INTO WAREHOUSE CSV DATA ('nodelocal://1/assets/processed/warehouse/warehouse.csv') WITH skip='1', nullif='NULL';
IMPORT INTO DISTRICT CSV DATA ('nodelocal://1/assets/processed/district/district.csv') WITH skip='1', nullif='NULL';
IMPORT INTO CUSTOMER CSV DATA ('nodelocal://1/assets/processed/customer/customer.csv') WITH skip='1', nullif='NULL';
IMPORT INTO ITEM CSV DATA ('nodelocal://1/assets/processed/item/item.csv') WITH skip='1', nullif='NULL';
IMPORT INTO STOCK CSV DATA ('nodelocal://1/assets/processed/stock/stock.csv') WITH skip='1', nullif='NULL';
IMPORT INTO ORDERS CSV DATA ('nodelocal://1/assets/processed/order/order.csv') WITH skip='1', nullif='NULL';
IMPORT INTO ORDER_LINE CSV DATA ('nodelocal://1/assets/processed/orderline/orderline.csv') WITH skip='1', nullif='NULL';
