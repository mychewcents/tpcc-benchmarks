package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
)

type sqlStatement struct {
	baseStatement string
	filePath      string
}

func exportCSV(c config.Configuration) {
	hostName := fmt.Sprintf("%s:%d", c.HostNode.Host, c.HostNode.Port)
	sqls := []sqlStatement{
		{
			baseStatement: "SELECT * FROM WAREHOUSE",
			filePath:      "assets/data/processed/warehouse/warehouse.csv",
		},
		{
			baseStatement: "SELECT * FROM DISTRICT",
			filePath:      "assets/data/processed/district/district.csv",
		},
		{
			baseStatement: "SELECT * FROM CUSTOMER",
			filePath:      "assets/data/processed/customer/customer.csv",
		},
		{
			baseStatement: "SELECT * FROM ITEM",
			filePath:      "assets/data/processed/item/item.csv",
		},
		{
			baseStatement: "SELECT * FROM STOCK",
			filePath:      "assets/data/processed/stock/stock.csv",
		},
		{
			baseStatement: "SELECT * FROM ORDERS",
			filePath:      "assets/data/processed/order/order.csv",
		},
		{
			baseStatement: "SELECT * FROM ORDER_LINE",
			filePath:      "assets/data/processed/orderline/orderline.csv",
		},
	}

	for _, value := range sqls {
		cmd := &exec.Cmd{
			Path: "scripts/export_data.sh",
			Args: []string{"scripts/export_data.sh",
				hostName,
				value.baseStatement,
				value.filePath,
			},
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Dir:    ".",
		}

		if err := cmd.Start(); err != nil {
			log.Fatalf("error occured in %s. Err: %v", value.filePath, err)
			return
		}
		if err := cmd.Wait(); err != nil {
			log.Fatalf("error occured in %s. Err: %v", value.filePath, err)
			return
		}
		log.Printf("Completed exporting: %s", value.filePath)
	}
	exportPartitionsCSV(c)

	log.Printf("Completed exporting the database")
}

func exportPartitionsCSV(c config.Configuration) {
	log.Printf("Starting the export of partitions...")

	sqls := []sqlStatement{
		{
			baseStatement: "SELECT * FROM ORDERS_WID_DID",
			filePath:      "order",
		},
		{
			baseStatement: "SELECT * FROM ORDER_LINE_WID_DID",
			filePath:      "orderline",
		},
		{
			baseStatement: "SELECT * FROM ORDER_ITEMS_CUSTOMERS_WID_DID",
			filePath:      "itempairs",
		},
	}

	ch := make(chan bool, 300)

	for _, value := range sqls {
		log.Printf("Starting the export of: %s", value.filePath)
		baseSQLStatement := value.baseStatement
		for w := 1; w <= 10; w++ {
			for d := 1; d <= 10; d++ {
				finalSQLStatement := strings.ReplaceAll(baseSQLStatement, "WID", fmt.Sprintf("%d", w))
				finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "DID", fmt.Sprintf("%d", d))
				fileName := fmt.Sprintf("assets/data/processed/%s/%d_%d.csv", value.filePath, w, d)

				go exportPartitionsCSVParallel(c, w, d, finalSQLStatement, fileName, value.filePath, ch)
			}
		}
	}

	totalCount := 0
	executed := false
	for i := 1; i <= 300; i++ {
		executed = <-ch
		if executed {
			totalCount++
		}
	}

	log.Printf("Completed the export of partitions. Total Count: %d ; Should be 300.", totalCount)
}

func exportPartitionsCSVParallel(c config.Configuration, w, d int, finalSQLStatement, fileName, filePath string, ch chan bool) {
	hostName := fmt.Sprintf("%s:%d", c.HostNode.Host, c.HostNode.Port)
	cmd := &exec.Cmd{
		Path: "scripts/export_data.sh",
		Args: []string{"scripts/export_data.sh",
			hostName,
			finalSQLStatement,
			fileName,
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Dir:    ".",
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("error occured in %s for %d %d. Err: %v", filePath, w, d, err)
		ch <- false
	}
	if err := cmd.Wait(); err != nil {
		log.Fatalf("error occured in %s for %d %d. Err: %v", filePath, w, d, err)
		ch <- false
	}

	log.Printf("Completed the partition for %s warehouse: %d", filePath, w)
	ch <- true
}
