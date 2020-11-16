package services

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// ExecuteSQLService interface to executing sqls
type ExecuteSQLService interface {
	Execute(sqlString string) error
	ExecutePartitions(warehouses, districts int, baseSQLStatement string) error
}

type executeSQLServiceImpl struct {
	db *sql.DB
}

// CreateExecuteSQLService creates a new service to execute SQL Scripts
func CreateExecuteSQLService(db *sql.DB) ExecuteSQLService {
	return &executeSQLServiceImpl{db: db}
}

func (eqs *executeSQLServiceImpl) Execute(sqlString string) error {
	for _, value := range strings.Split(sqlString, ";") {
		if _, err := eqs.db.Exec(value); err != nil {
			return fmt.Errorf("error occurred in executing the query: %s. \nErr: %v", value, err)
		}
	}

	return nil
}

func (eqs *executeSQLServiceImpl) ExecutePartitions(warehouses, districts int, baseSQLStatement string) error {
	for w := 1; w <= warehouses; w++ {
		for d := 1; d <= districts; d++ {
			finalSQLStatement := strings.ReplaceAll(baseSQLStatement, "WID", strconv.Itoa(w))
			finalSQLStatement = strings.ReplaceAll(finalSQLStatement, "DID", strconv.Itoa(d))

			if err := eqs.Execute(finalSQLStatement); err != nil {
				return err
			}
		}
	}

	return nil
}
