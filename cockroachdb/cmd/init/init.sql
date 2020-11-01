DROP TABLE IF EXISTS STOCK;
DROP TABLE IF EXISTS ITEM;
DROP TABLE IF EXISTS ORDER_LINE;
DROP TABLE IF EXISTS ORDERS;
DROP TABLE IF EXISTS CUSTOMER;
DROP TABLE IF EXISTS DISTRICT;
DROP TABLE IF EXISTS WAREHOUSE;

CREATE TABLE IF NOT EXISTS WAREHOUSE (
  W_ID 						INT PRIMARY KEY,
  W_NAME 					STRING NOT NULL,
  W_STREET_1 			STRING NOT NULL,
  W_STREET_2 			STRING NOT NULL,
  W_CITY 					STRING NOT NULL,
  W_STATE 				STRING NOT NULL,
  W_ZIP						STRING NOT NULL,
  W_TAX						DECIMAL(4,4) NOT NULL,
  W_YTD						DECIMAL(12,2) NOT NULL	
);

IMPORT INTO WAREHOUSE (
  W_ID,
  W_NAME,
  W_STREET_1,
  W_STREET_2,
  W_CITY,
  W_STATE,
  W_ZIP,
  W_TAX,
  W_YTD
) CSV DATA ('nodelocal://1/assets/raw/warehouse.csv');

DROP TABLE IF EXISTS DISTRICT;
CREATE TABLE IF NOT EXISTS DISTRICT (
  D_W_ID          INT REFERENCES WAREHOUSE(W_ID),
  D_ID            INT,
  D_NAME          STRING NOT NULL,
  D_STREET_1      STRING NOT NULL,
  D_STREET_2      STRING NOT NULL,
  D_CITY          STRING NOT NULL,
  D_STATE         STRING NOT NULL,
  D_ZIP           STRING NOT NULL,
  D_TAX           DECIMAL(4,4) NOT NULL,
  D_YTD           DECIMAL(12,2) NOT NULL,
  D_NEXT_O_ID     INT NOT NULL,
  INDEX (D_W_ID),
  PRIMARY KEY (D_W_ID, D_ID)
);

IMPORT INTO DISTRICT (
  D_W_ID,
  D_ID,
  D_NAME,
  D_STREET_1,
  D_STREET_2,
  D_CITY,
  D_STATE,
  D_ZIP,
  D_TAX,
  D_YTD,
  D_NEXT_O_ID
) CSV DATA ('nodelocal://1/assets/raw/district.csv');

ALTER TABLE DISTRICT ADD COLUMN D_W_TAX DECIMAL(4,4);
UPDATE DISTRICT SET D_W_TAX = W_TAX FROM WAREHOUSE WHERE DISTRICT.D_W_ID = WAREHOUSE.W_ID;

DROP TABLE IF EXISTS CUSTOMER;
CREATE TABLE IF NOT EXISTS CUSTOMER (
  C_W_ID                INT,
  C_D_ID                INT,
  C_ID                  INT,
  C_FIRST               STRING NOT NULL,
  C_MIDDLE              STRING NOT NULL,
  C_LAST                STRING NOT NULL,
  C_STREET_1            STRING NOT NULL,
  C_STREET_2            STRING NOT NULL,
  C_CITY                STRING NOT NULL,
  C_STATE               STRING NOT NULL,
  C_ZIP                 STRING NOT NULL,
  C_PHONE               STRING NOT NULL,
  C_SINCE               TIMESTAMP,
  C_CREDIT              STRING NOT NULL,
  C_CREDIT_LIM          DECIMAL(12,2),
  C_DISCOUNT            DECIMAL(4,4),
  C_BALANCE             DECIMAL(12,2),
  C_YTD_PAYMENT         FLOAT,
  C_PAYMENT_CNT         INT,
  C_DELIVERY_CNT        INT,
  C_DATA                STRING NOT NULL,
  INDEX (C_W_ID, C_D_ID),
  INDEX (C_BALANCE),
  PRIMARY KEY (C_W_ID, C_D_ID, C_ID),
  CONSTRAINT FK_CUSTOMERS FOREIGN KEY (C_W_ID, C_D_ID) REFERENCES DISTRICT (D_W_ID, D_ID) 
);

IMPORT INTO CUSTOMER (
  C_W_ID,
  C_D_ID,
  C_ID,
  C_FIRST,
  C_MIDDLE,
  C_LAST,
  C_STREET_1,
  C_STREET_2,
  C_CITY,
  C_STATE,
  C_ZIP,
  C_PHONE,
  C_SINCE,
  C_CREDIT,
  C_CREDIT_LIM,
  C_DISCOUNT,
  C_BALANCE,
  C_YTD_PAYMENT,
  C_PAYMENT_CNT,
  C_DELIVERY_CNT,
  C_DATA
) CSV DATA ('nodelocal://1/assets/raw/customer.csv');

DROP TABLE IF EXISTS ORDERS;
CREATE TABLE IF NOT EXISTS ORDERS (
  O_W_ID int,
  O_D_ID int,
  O_ID int,
  O_C_ID int NULL,
  O_CARRIER_ID int,
  O_OL_CNT decimal(2,0),
  O_ALL_LOCAL DECIMAL(1,0),
  O_ENTRY_D timestamp,
  INDEX (O_C_ID),
  PRIMARY KEY (O_W_ID, O_D_ID, O_ID),
  CONSTRAINT FK_ORDERS FOREIGN KEY (O_W_ID, O_D_ID, O_C_ID) REFERENCES CUSTOMER (C_W_ID, C_D_ID, C_ID)
);

IMPORT INTO ORDERS (
  O_W_ID,
  O_D_ID,
  O_ID,
  O_C_ID,
  O_CARRIER_ID,
  O_OL_CNT,
  O_ALL_LOCAL,
  O_ENTRY_D
) CSV DATA ('nodelocal://1/assets/raw/order.csv') WITH nullif='null';

ALTER TABLE ORDERS ADD COLUMN O_TOTAL_AMOUNT DECIMAL(12, 2) DEFAULT 0.0;
ALTER TABLE ORDERS ADD COLUMN O_DELIVERY_D TIMESTAMP;
ALTER TABLE ORDERS ADD COLUMN O_OL_ITEM_IDS STRING;

DROP TABLE IF EXISTS ORDER_LINE;
CREATE TABLE IF NOT EXISTS ORDER_LINE (
  OL_W_ID int,
  OL_D_ID int,
  OL_O_ID int,
  OL_NUMBER int,
  OL_I_ID int,
  OL_DELIVERY_D timestamp,
  OL_AMOUNT decimal(6,2),
  OL_SUPPLY_W_ID int,
  OL_QUANTITY decimal(2,0),
  OL_DIST_INFO char(24),
  INDEX (OL_O_ID),
  INDEX (OL_I_ID),
  PRIMARY KEY (OL_W_ID, OL_D_ID, OL_O_ID, OL_NUMBER),
  CONSTRAINT FK_ORDER_LINE FOREIGN KEY (OL_W_ID, OL_D_ID, OL_O_ID) REFERENCES ORDERS (O_W_ID, O_D_ID, O_ID)
);

IMPORT INTO ORDER_LINE (
  OL_W_ID,
  OL_D_ID,
  OL_O_ID,
  OL_NUMBER,
  OL_I_ID,
  OL_DELIVERY_D,
  OL_AMOUNT,
  OL_SUPPLY_W_ID,
  OL_QUANTITY,
  OL_DIST_INFO
) CSV DATA ('nodelocal://1/assets/raw/order-line.csv') WITH nullif='null';

DROP TABLE IF EXISTS ITEM;
CREATE TABLE IF NOT EXISTS ITEM (
  I_ID int PRIMARY KEY,
  I_NAME varchar(24),
  I_PRICE decimal(5,2),
  I_IM_ID int,
  I_DATA varchar(50)
);

IMPORT INTO ITEM (
  I_ID,
  I_NAME,
  I_PRICE,
  I_IM_ID,
  I_DATA
) CSV DATA ('nodelocal://1/assets/raw/item.csv');

DROP TABLE IF EXISTS STOCK;
CREATE TABLE IF NOT EXISTS STOCK (
  S_W_ID int,
  S_I_ID int,
  S_QUANTITY decimal(4,0),
  S_YTD decimal(8,2),
  S_ORDER_CNT int,
  S_REMOTE_CNT int,
  S_DIST_01 char(24),
  S_DIST_02 char(24),
  S_DIST_03 char(24),
  S_DIST_04 char(24),
  S_DIST_05 char(24),
  S_DIST_06 char(24),
  S_DIST_07 char(24),
  S_DIST_08 char(24),
  S_DIST_09 char(24),
  S_DIST_10 char(24),
  S_DATA varchar(50),
  INDEX (S_W_ID),
  PRIMARY KEY (S_W_ID, S_I_ID),
  CONSTRAINT FK_STOCK_WAREHOUSE FOREIGN KEY (S_W_ID) REFERENCES WAREHOUSE (W_ID),
  CONSTRAINT FK_STOCK_ITEM FOREIGN KEY (S_I_ID) REFERENCES ITEM (I_ID)
);

IMPORT INTO STOCK (
  S_W_ID,
  S_I_ID,
  S_QUANTITY,
  S_YTD,
  S_ORDER_CNT,
  S_REMOTE_CNT,
  S_DIST_01,
  S_DIST_02,
  S_DIST_03,
  S_DIST_04,
  S_DIST_05,
  S_DIST_06,
  S_DIST_07,
  S_DIST_08,
  S_DIST_09,
  S_DIST_10,
  S_DATA
) CSV DATA ('nodelocal://1/assets/raw/stock.csv');

ALTER TABLE STOCK ADD COLUMN S_I_NAME STRING;
ALTER TABLE STOCK ADD COLUMN S_I_PRICE DECIMAL(5, 2);

UPDATE STOCK SET S_I_NAME = I_NAME, S_I_PRICE = I_PRICE FROM ITEM WHERE STOCK.S_I_ID = ITEM.I_ID;
