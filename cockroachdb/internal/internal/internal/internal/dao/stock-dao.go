package dao

import (
	"database/sql"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dbdatamodel"
)

// StockDao creates the Dao object for the Stock table
type StockDao interface {
	GetStockDetails(warehouseID, itemID int) dbdatamodel.Stock
}

type stockDaoImpl struct {
	db *sql.DB
}

// GetStockDao creates the new StockDao object
func GetStockDao(db *sql.DB) StockDao {
	return &stockDaoImpl{db: db}
}

// GetStockDetails gets the stock details
func (sd *stockDaoImpl) GetStockDetails(warehouseID, itemID int) dbdatamodel.Stock {
	
}
