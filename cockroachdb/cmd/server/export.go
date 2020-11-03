package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
)

func export(c config.Configuration) {
	hostName := fmt.Sprintf("%s:%d", c.HostNode.Host, c.HostNode.Port)
	baseSQLStatement := fmt.Sprintf("SELECT * FROM ORDER_ITEMS_CUSTOMERS_WID_DID;")
	for w := 1; w <= 10; w++ {
		for d := 1; d <= 10; d++ {
			finalSQLStatement := strings.ReplaceAll(baseSQLStatement, "WID", fmt.Sprintf("%d", w))
			finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "DID", fmt.Sprintf("%d", d))
			fileName := fmt.Sprintf("assets/processed/item_pairs/%d_%d.csv", w, d)

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
