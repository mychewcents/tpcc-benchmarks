package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/cdbconn"
)

var db *sql.DB

func init() {
	var err error
	db, err = cdbconn.CreateConnection()
	if err != nil {
		panic(err)
	}
}

func main() {
	createOrdersTables(2, 2)
	createOrderLinesTables()
}

func createOrdersTables(warehouses, districts int) {
	baseSQLStatement := `
		DROP TABLE IF EXISTS defaultdb.ORDERS_WID_DID;
		
		CREATE TABLE IF NOT EXISTS defaultdb.ORDERS_WID_DID (
			O_W_ID int,
			O_D_ID int,
			O_ID int,
			O_C_ID int NULL,
			O_CARRIER_ID int DEFAULT NULL,
			O_OL_CNT decimal(2,0),
			O_ALL_LOCAL DECIMAL(1,0),
			O_ENTRY_D timestamp DEFAULT CURRENT_TIMESTAMP,
			O_TOTAL_AMOUNT decimal(12,2),
			INDEX (O_C_ID),
			PRIMARY KEY (O_W_ID, O_D_ID, O_ID),
			CONSTRAINT FK_ORDERS FOREIGN KEY (O_W_ID, O_D_ID, O_C_ID) REFERENCES defaultdb.CUSTOMER (C_W_ID, C_D_ID, C_ID)
		);
		
		INSERT INTO defaultdb.ORDERS_WID_DID (O_W_ID, O_D_ID, O_ID, O_C_ID, O_CARRIER_ID, O_OL_CNT, O_ALL_LOCAL, O_ENTRY_D, O_TOTAL_AMOUNT) 
		SELECT O_W_ID, O_D_ID, O_ID, O_C_ID, O_CARRIER_ID, O_OL_CNT, O_ALL_LOCAL, O_ENTRY_D, 0.0 FROM defaultdb.ORDERS 
		WHERE O_W_ID = WID AND O_D_ID = DID;
	`

	for i := 1; i <= warehouses; i++ {
		for j := 1; j <= districts; j++ {
			finalSQLStatement := strings.ReplaceAll(baseSQLStatement, "WID", strconv.Itoa(i))
			finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "DID", strconv.Itoa(j))

			// fmt.Println(finalSQLStatement)
			_, err := db.Exec(finalSQLStatement)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("Compelete: %d %d", i, j)
		}
	}
}

func createOrderLinesTables() {

}
