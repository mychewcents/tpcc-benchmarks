package dao

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dbdatamodel"
)

// WarehouseDao interface to the Warehouse table operations
type WarehouseDao interface {
	GetAddress(tx *sql.Tx, warehouseID int) (*dbdatamodel.Address, error)
	GetWarehouseNames(warehouseIDs []int) (map[int]string, error)
}

type warehouseDaoImpl struct {
	db *sql.DB
}

// CreateWarehouseDao creates new object for the WarehouseDao
func CreateWarehouseDao(db *sql.DB) WarehouseDao {
	return &warehouseDaoImpl{db: db}
}

// GetAddress returns the address of the warehouse id
func (wd *warehouseDaoImpl) GetAddress(tx *sql.Tx, warehouseID int) (addr *dbdatamodel.Address, err error) {

	sqlStatement := fmt.Sprintf("SELECT W_STREET_1, W_STREET_2, W_CITY, W_STATE, W_ZIP FROM WAREHOUSE WHERE W_ID = %d", warehouseID)

	if err = tx.QueryRow(sqlStatement).Scan(&addr.Street1, &addr.Street2, &addr.City, &addr.State, &addr.Zip); err != nil {
		return nil, fmt.Errorf("error occurred in reading the warehouse details. Err: %v", err)
	}

	return
}

func (wd *warehouseDaoImpl) GetWarehouseNames(warehouseIDs []int) (result map[int]string, err error) {
	var whereClauseBuilder strings.Builder

	for _, value := range warehouseIDs {
		whereClauseBuilder.WriteString(fmt.Sprintf("%d,", value))
	}

	whereClauseString := whereClauseBuilder.String()
	whereClauseString = whereClauseString[:len(whereClauseString)-1]

	sqlStatement := fmt.Sprintf("SELECT W_ID, W_NAME FROM Warehouse WHERE W_ID IN (%s)", whereClauseString)

	rows, err := wd.db.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("error occurred in getting the warehouse names. Err: %v", err)
	}
	defer rows.Close()

	var wID int
	var name string

	for rows.Next() {
		if err := rows.Scan(&wID, &name); err != nil {
			return nil, fmt.Errorf("error occurred in scanning the warehouse details. Err: %v", err)
		}

		result[wID] = name
	}

	return
}
