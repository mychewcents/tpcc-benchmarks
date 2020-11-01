package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/cdbconn"
)

var db *sql.DB

func init() {
	var err error

	if len(os.Args) != 2 {
		panic("Missing configuration file path")
	}
	db, err = cdbconn.CreateConnection(os.Args[1])
	if err != nil {
		panic(err)
	}

	fileName := fmt.Sprintf("logs/logs_init_%s", time.Now())
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		return
	}

	log.SetOutput(file)
}

func main() {
	log.Println("Starting the loading of the raw database files")
	if err := loadRawDataset(db, "cmd/init/init.sql"); err != nil {
		fmt.Println(err)
		return
	}
	log.Println("Finished the loading of the raw database files")

	log.Println("Starting the creation of partitioned Orders table...")
	if err := createOrdersTables(10, 10); err != nil {
		fmt.Println(err)
		return
	}
	log.Println("Finished the creation of partitionied Orders table...")

	log.Println("Starting the creation of partitioned Order Line table...")
	if err := createOrderLinesTables(10, 10); err != nil {
		fmt.Println(err)
		return
	}
	log.Println("Finished the creation of partitionied Order Line table...")

	log.Println("Updating the Orders table with the Total Amount for each order...")
	if err := updateOrdersTotalAmount(10, 10); err != nil {
		fmt.Println(err)
		return
	}
	log.Println("Finished updated the Orders table with the Total Amount...")

	log.Println("Starting the creation of Customer <-> Two Item Pair table...")
	if err := createOrderItemsPairTables(10); err != nil {
		fmt.Println(err)
		return
	}
	log.Println("Finished the creation of ORDER_ITEMS_CUSTOMERS tables for each warehouse")

	fmt.Println("Done")
}

func loadRawDataset(db *sql.DB, file string) error {
	initSQL, err := os.Open(file)
	if err != nil {
		log.Fatalf("Err: %v", err)
		return errors.New("error occurred. Please check the logs")
	}
	defer initSQL.Close()

	byteValue, _ := ioutil.ReadAll(initSQL)

	var finalQueryBuilder strings.Builder
	finalQueryBuilder.WriteString(string(byteValue))

	for _, value := range strings.Split(finalQueryBuilder.String(), ";") {
		log.Println(value)

		if _, err = db.Exec(value); err != nil {
			log.Fatalf("Err: %v", err)
			return errors.New("error occurred. Please check the logs")
		}
	}
	return nil
}

func createOrdersTables(warehouses, districts int) error {
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
			O_DELIVERY_D timestamp DEFAULT NULL,
			O_OL_ITEM_IDS string DEFAULT NULL,
			INDEX (O_C_ID, O_ID DESC),
			INDEX (O_CARRIER_ID, O_ID),
			PRIMARY KEY (O_W_ID, O_D_ID, O_ID),
			CONSTRAINT FK_ORDERS FOREIGN KEY (O_W_ID, O_D_ID, O_C_ID) REFERENCES defaultdb.CUSTOMER (C_W_ID, C_D_ID, C_ID)
		);
		
		INSERT INTO defaultdb.ORDERS_WID_DID (O_W_ID, O_D_ID, O_ID, O_C_ID, O_CARRIER_ID, O_OL_CNT, O_ALL_LOCAL, O_ENTRY_D, O_TOTAL_AMOUNT) 
		SELECT O_W_ID, O_D_ID, O_ID, O_C_ID, CASE WHEN O_CARRIER_ID IS NULL THEN 0 ELSE O_CARRIER_ID END, O_OL_CNT, O_ALL_LOCAL, O_ENTRY_D, 0.0 FROM defaultdb.ORDERS 
		WHERE O_W_ID = WID AND O_D_ID = DID;
	`

	errFound := false
	for i := 1; i <= warehouses; i++ {
		for j := 1; j <= districts; j++ {
			finalSQLStatement := strings.ReplaceAll(baseSQLStatement, "WID", strconv.Itoa(i))
			finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "DID", strconv.Itoa(j))

			log.Println(finalSQLStatement)
			if _, err := db.Exec(finalSQLStatement); err != nil {
				log.Fatalf("Err: %v", err)
				errFound = true
			}
		}
	}

	if errFound {
		return errors.New("error was found. Please check the logs")
	}
	return nil
}

func createOrderLinesTables(warehouses, districts int) error {
	baseSQLStatement := `
		DROP TABLE IF EXISTS defaultdb.ORDER_LINE_WID_DID;
		CREATE TABLE IF NOT EXISTS defaultdb.ORDER_LINE_WID_DID (
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
			PRIMARY KEY (OL_W_ID, OL_D_ID, OL_O_ID, OL_NUMBER)
			CONSTRAINT FK_ORDER_LINE FOREIGN KEY (OL_W_ID, OL_D_ID, OL_O_ID) REFERENCES defaultdb.ORDERS_WID_DID (O_W_ID, O_D_ID, O_ID)
		);
		
		INSERT INTO defaultdb.ORDER_LINE_WID_DID 
		SELECT * FROM defaultdb.ORDER_LINE
		WHERE OL_W_ID = WID AND OL_D_ID = DID;
	`

	errFound := false
	for i := 1; i <= warehouses; i++ {
		for j := 1; j <= districts; j++ {
			finalSQLStatement := strings.ReplaceAll(baseSQLStatement, "WID", strconv.Itoa(i))
			finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "DID", strconv.Itoa(j))

			log.Println(finalSQLStatement)
			if _, err := db.Exec(finalSQLStatement); err != nil {
				log.Fatalf("Err: %v", err)
				errFound = true
			}
		}
	}

	if errFound {
		return errors.New("error was found. Please check the logs")
	}
	return nil
}

func updateOrdersTotalAmount(warehouses, districts int) error {
	baseSQLStatement := "UPDATE ORDERS_WID_DID SET O_TOTAL_AMOUNT = (SELECT SUM(OL_AMOUNT) FROM ORDER_LINE_WID_DID WHERE OL_O_ID = O_ID)"

	errFound := false
	for i := 1; i <= warehouses; i++ {
		for j := 1; j <= districts; j++ {
			finalSQLStatement := strings.ReplaceAll(baseSQLStatement, "WID", strconv.Itoa(i))
			finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "DID", strconv.Itoa(j))

			log.Println(finalSQLStatement)
			if _, err := db.Exec(finalSQLStatement); err != nil {
				log.Fatalf("Err: %v", err)
				errFound = true
			}
		}
	}

	if errFound {
		return errors.New("error was found. Please check the logs")
	}
	return nil
}

func createOrderItemsPairTables(warehouses int) error {
	baseSQLStatement := `
		DROP TABLE IF EXISTS defaultdb.ORDER_ITEMS_CUSTOMERS_WID;
		
		CREATE TABLE IF NOT EXISTS defaultdb.ORDER_ITEMS_CUSTOMERS_WID (
			IC_W_ID int,
			IC_D_ID int,
			IC_C_ID int,
			IC_I_1_ID int,
			IC_I_2_ID int,
			INDEX (IC_W_ID, IC_D_ID, IC_C_ID),
			PRIMARY KEY (IC_W_ID, IC_D_ID, IC_C_ID, IC_I_1_ID, IC_I_2_ID),
			CONSTRAINT FK_ORDERS FOREIGN KEY (IC_W_ID, IC_D_ID, IC_C_ID) REFERENCES defaultdb.CUSTOMER (C_W_ID, C_D_ID, C_ID)
		);
	`

	errFound := false
	for i := 1; i <= warehouses; i++ {
		finalSQLStatement := strings.ReplaceAll(baseSQLStatement, "WID", strconv.Itoa(i))

		log.Println(finalSQLStatement)
		if _, err := db.Exec(finalSQLStatement); err != nil {
			log.Fatalf("Err: %v", err)
			errFound = true
		}
	}

	if errFound {
		return errors.New("error was found. Please check the logs")
	}
	return nil
}
