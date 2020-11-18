package services

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/server"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/common/config"
)

// ExportTablesService interface to the service to export tables
type ExportTablesService interface {
	Export() error
}

type exportTablesServiceImpl struct {
	c  config.Configuration
	db *sql.DB
}

// CreateExportTablesService creates a new service to export the tables
func CreateExportTablesService(c config.Configuration, db *sql.DB) ExportTablesService {
	return &exportTablesServiceImpl{c: c, db: db}
}

func (ets *exportTablesServiceImpl) Export() (err error) {
	if err = ets.exportParentTables(); err != nil {
		return err
	}

	if err = ets.exportPartitionedTables(); err != nil {
		return err
	}

	return nil
}

func (ets *exportTablesServiceImpl) exportParentTables() (err error) {
	hostName := fmt.Sprintf("%s:%d", ets.c.HostNode.Host, ets.c.HostNode.Port)
	sqls := []*models.TableExportParams{
		{
			BaseSQLStatement: "SELECT * FROM WAREHOUSE",
			ExportPath:       "assets/data/processed/warehouse/warehouse.csv",
		},
		{
			BaseSQLStatement: "SELECT * FROM DISTRICT",
			ExportPath:       "assets/data/processed/district/district.csv",
		},
		{
			BaseSQLStatement: "SELECT * FROM CUSTOMER",
			ExportPath:       "assets/data/processed/customer/customer.csv",
		},
		{
			BaseSQLStatement: "SELECT * FROM ITEM",
			ExportPath:       "assets/data/processed/item/item.csv",
		},
		{
			BaseSQLStatement: "SELECT * FROM STOCK",
			ExportPath:       "assets/data/processed/stock/stock.csv",
		},
		{
			BaseSQLStatement: "SELECT * FROM ORDERS",
			ExportPath:       "assets/data/processed/order/order.csv",
		},
		{
			BaseSQLStatement: "SELECT * FROM ORDER_LINE",
			ExportPath:       "assets/data/processed/orderline/orderline.csv",
		},
	}

	for _, value := range sqls {
		cliArgs := []string{"scripts/export_data.sh",
			hostName,
			value.BaseSQLStatement,
			value.ExportPath,
		}

		if err = server.Execute(cliArgs); err != nil {
			return err
		}

		log.Printf("Completed exporting: %s", value.ExportPath)
	}

	return nil
}

func (ets *exportTablesServiceImpl) exportPartitionedTables() (err error) {
	hostName := fmt.Sprintf("%s:%d", ets.c.HostNode.Host, ets.c.HostNode.Port)
	sqls := []*models.TableExportParams{
		{
			BaseSQLStatement: "SELECT * FROM ORDERS_WID_DID",
			ExportPath:       "assets/data/processed/order/WID_DID.csv",
		},
		{
			BaseSQLStatement: "SELECT * FROM ORDER_LINE_WID_DID",
			ExportPath:       "assets/data/processed/orderline/WID_DID.csv",
		},
		{
			BaseSQLStatement: "SELECT * FROM ORDER_ITEMS_CUSTOMERS_WID_DID",
			ExportPath:       "assets/data/processed/itempairs/WID_DID.csv",
		},
	}

	for _, value := range sqls {
		baseSQLStatement := value.BaseSQLStatement
		baseExportPath := value.ExportPath

		for w := 1; w <= 10; w++ {
			for d := 1; d <= 10; d++ {
				finalSQLStatement := strings.ReplaceAll(baseSQLStatement, "WID", fmt.Sprintf("%d", w))
				finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "DID", fmt.Sprintf("%d", d))
				finalExportPath := strings.ReplaceAll(baseExportPath, "WID", fmt.Sprintf("%d", w))
				finalExportPath = strings.ReplaceAll(finalExportPath, "DID", fmt.Sprintf("%d", d))

				cliArgs := []string{"scripts/export_data.sh",
					hostName,
					finalSQLStatement,
					finalExportPath,
				}

				if err = server.Execute(cliArgs); err != nil {
					return err
				}
				log.Printf("Completed exporting: %s", value.ExportPath)
			}
		}
	}

	return nil
}
