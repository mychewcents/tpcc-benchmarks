package dao

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dbdatamodel"
)

// DistrictDao interface to the functions accessing district table
type DistrictDao interface {
	GetNewOrderIDAndTaxRates(warehouseID, districtID int) (int, float64, float64, error)
	GetLastOrderID(warehouseID, districtID int) (int, error)
	AddPaymentToDistrict(tx *sql.Tx, warehouseID, districtID int, amount float64) (*dbdatamodel.Address, error)
	GetDistrictNames(warehouseIDs, districtIDs []int) (map[int]map[int]string, error)
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

func (dd *districtDaoImpl) AddPaymentToDistrict(tx *sql.Tx, warehouseID, districtID int, amount float64) (addr *dbdatamodel.Address, err error) {
	sqlStatement := fmt.Sprintf(`
		UPDATE DISTRICT SET 
		D_YTD = D_YTD + %f 
		WHERE (D_W_ID, D_ID) = (%d, %d)
		RETURNING D_STREET_1, D_STREET_2, D_CITY, D_STATE, D_ZIP`,
		amount, warehouseID, districtID)

	if err := tx.QueryRow(sqlStatement).Scan(&addr.Street1, &addr.Street2, &addr.City, &addr.State, &addr.Zip); err != nil {
		return nil, fmt.Errorf("error occurred in updating the district table. Err: %v", err)
	}

	return
}

func (dd *districtDaoImpl) GetDistrictNames(warehouseIDs, districtIDs []int) (result map[int]map[int]string, err error) {

	var whereClauseBuilder strings.Builder

	for _, wValue := range warehouseIDs {
		for _, dValue := range districtIDs {
			whereClauseBuilder.WriteString(fmt.Sprintf("(%d, %d),", wValue, dValue))
		}
	}

	whereClauseString := whereClauseBuilder.String()
	whereClauseString = whereClauseString[:len(whereClauseString)-1]

	sqlStatement := fmt.Sprintf("SELECT D_W_ID, D_ID, D_NAME FROM District WHERE (D_W_ID, D_ID) IN (%s)", whereClauseString)

	rows, err := dd.db.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("error occurred in reading the names of the districts. Err: %v", err)
	}
	defer rows.Close()

	var wID, dID int
	var name string

	for rows.Next() {
		if err := rows.Scan(&wID, &dID, &name); err != nil {
			return nil, fmt.Errorf("error occurred in scanning the rows from the result set. Err: %v", err)
		}
		if result[wID] == nil {
			result[wID] = make(map[int]string)
		}
		result[wID][dID] = name
	}

	return
}
