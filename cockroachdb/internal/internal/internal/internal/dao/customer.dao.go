package dao

import (
	"database/sql"
	"fmt"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dbdatamodel"
)

// CustomerDao represents the interface to get the Customer details from the DB
type CustomerDao interface {
	GetCustomerDetails(warehouseID, districtID, customerID int) (customer *dbdatamodel.Customer, err error)
}

type customerDaoImpl struct {
	db *sql.DB
}

// CreateCustomerDao creates new customer dao
func CreateCustomerDao(db *sql.DB) CustomerDao {
	return &customerDaoImpl{db: db}
}

// GetCustomerDetails gets the customer details from the database
func (cs *customerDaoImpl) GetCustomerDetails(warehouseID, districtID, customerID int) (customer *dbdatamodel.Customer, err error) {
	sqlStatement := fmt.Sprintf("SELECT C_LAST, C_CREDIT, C_DISCOUNT FROM CUSTOMER WHERE C_W_ID = $1 AND C_D_ID = $2 AND C_ID = $3")

	var lastName, credit string
	var discount float64

	row := cs.db.QueryRow(sqlStatement, warehouseID, districtID, customerID)
	if err := row.Scan(&lastName, &credit, &discount); err != nil {
		return nil, fmt.Errorf("error occured in getting the customer details. Err: %v", err)
	}

	customer = &dbdatamodel.Customer{
		WarehouseID: warehouseID,
		DistrictID:  districtID,
		CustomerID:  customerID,
		LastName:    lastName,
		Credit:      credit,
		Discount:    discount,
	}

	return
}
