package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
)

type sqlStatement struct {
	baseStatement string
	filePath      string
}

func export(c config.Configuration) {
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
		// {
		// 	baseStatement: "SELECT * FROM ORDER_ITEMS_CUSTOMERS_WID_DID",
		// 	filePath:      "itempairs",
		// },
	}

	for _, value := range sqls {
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

				cmd.Start()
				cmd.Wait()
			}
			fmt.Printf("Completed: w = %d\n", w)
		}
	}
}
