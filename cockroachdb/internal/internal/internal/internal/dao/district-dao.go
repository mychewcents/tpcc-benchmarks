package dao

import (
	"database/sql"
	"fmt"
)

// DistrictDao interface to the functions accessing district table
type DistrictDao interface {
	GetNewOrderIDAndTaxRates(warehouseID, districtID int) (int, float64, float64, error)
	GetLastOrderID(warehouseID, districtID int) (int, error)
}

type districtDaoImpl struct {
	db *sql.DB
}

// CreateDistrictDao creates new District Dao object
func CreateDistrictDao(db *sql.DB) DistrictDao {
	return &districtDaoImpl{db: db}
}

func (dd *districtDaoImpl) GetNewOrderIDAndTaxRates(warehouseID, districtID int) (newOrderID int, wTax, dTax float64, err error) {
	sqlStatement := fmt.Sprintf("UPDATE District SET D_NEXT_O_ID = D_NEXT_O_ID + 1 WHERE D_W_ID = $1 AND D_ID = $2 RETURNING D_NEXT_O_ID, D_TAX, D_W_TAX")

	row := dd.db.QueryRow(sqlStatement, warehouseID, districtID)
	if err := row.Scan(&newOrderID, &dTax, &wTax); err != nil {
		return 0, 0.0, 0.0, fmt.Errorf("error occured in updating the district table for the next order id. Err: %v", err)
	}

	return
}

func (dd *districtDaoImpl) GetLastOrderID(warehouseID, districtID int) (lastOrderID int, err error) {
	row := dd.db.QueryRow("SELECT d_next_o_id FROM district WHERE d_w_id=$1 AND d_id=$2", warehouseID, districtID)

	if err := row.Scan(&lastOrderID); err != nil {
		return lastOrderID, fmt.Errorf("error occurred in getting the next order id for the district. Err: %v", err)
	}

	return
}
