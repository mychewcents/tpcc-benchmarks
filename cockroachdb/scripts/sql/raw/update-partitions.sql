ALTER TABLE defaultdb.ORDER_LINE_WID_DID 
ADD CONSTRAINT FK_ORDER_LINE_WID_DID FOREIGN KEY (OL_W_ID, OL_D_ID, OL_O_ID) 
REFERENCES defaultdb.ORDERS_WID_DID (O_W_ID, O_D_ID, O_ID);

UPDATE ORDERS_WID_DID 
SET O_TOTAL_AMOUNT = (SELECT SUM(OL_AMOUNT) FROM ORDER_LINE_WID_DID WHERE OL_O_ID = O_ID);
