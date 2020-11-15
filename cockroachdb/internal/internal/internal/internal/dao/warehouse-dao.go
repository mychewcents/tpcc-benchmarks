package dao

import (
	"database/sql"
	"fmt"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dbdatamodel"
)

// WarehouseDao interface to the Warehouse table operations
type WarehouseDao interface {
	GetAddress(tx *sql.Tx, warehouseID int) (*dbdatamodel.Address, error)
}

type warehouseDaoImpl struct {
	db *sql.DB
}

// CreateWarhouseDao creates new object for the WarehouseDao
func CreateWarhouseDao(db *sql.DB) WarehouseDao {
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
