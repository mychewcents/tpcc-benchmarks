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

	hostName := fmt.Sprintf("%s:%d", c.HostNode.Host, c.HostNode.Port)
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

	for _, value := range sqls {
		log.Printf("Starting the export of: %s", value.filePath)
		baseSQLStatement := value.baseStatement
		for w := 1; w <= 10; w++ {
			for d := 1; d <= 10; d++ {
				finalSQLStatement := strings.ReplaceAll(baseSQLStatement, "WID", fmt.Sprintf("%d", w))
				finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "DID", fmt.Sprintf("%d", d))
				fileName := fmt.Sprintf("assets/data/processed/%s/%d_%d.csv", value.filePath, w, d)

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
					log.Fatalf("error occured in %s for %d %d. Err: %v", value.filePath, w, d, err)
					return
				}
				if err := cmd.Wait(); err != nil {
					log.Fatalf("error occured in %s for %d %d. Err: %v", value.filePath, w, d, err)
					return
				}
			}
			log.Printf("Completed the partition for warehouse: %d", w)
		}
		log.Printf("Completed the export of: %s", value.filePath)
	}
	log.Printf("Completed the export of partitions...")
}
