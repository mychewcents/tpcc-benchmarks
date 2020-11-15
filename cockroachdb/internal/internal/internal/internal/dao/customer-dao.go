package dao

import (
	"database/sql"
	"fmt"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dbdatamodel"
)

// CustomerDao represents the interface to get the Customer details from the DB
type CustomerDao interface {
	GetDetails(warehouseID, districtID, customerID int) (*dbdatamodel.Customer, error)
	UpdatePaymentDetails(tx *sql.Tx, warehouseID, districtID, customerID int, amount float64) (*dbdatamodel.Customer, error)
}

type customerDaoImpl struct {
	db *sql.DB
}

// CreateCustomerDao creates new customer dao
func CreateCustomerDao(db *sql.DB) CustomerDao {
	return &customerDaoImpl{db: db}
}

// GetDetails gets the customer details from the database
func (cs *customerDaoImpl) GetDetails(warehouseID, districtID, customerID int) (customer *dbdatamodel.Customer, err error) {
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

// UpdatePaymentDetails updates the new payment made by the customer
func (cs *customerDaoImpl) UpdatePaymentDetails(tx *sql.Tx, warehouseID, districtID, customerID int, amount float64) (customer *dbdatamodel.Customer, err error) {

	sqlStatement := fmt.Sprintf(`
		UPDATE CUSTOMER SET 
		(C_BALANCE, C_YTD_PAYMENT, C_PAYMENT_CNT) = (C_BALANCE - %f, C_YTD_PAYMENT + %[1]f, C_PAYMENT_CNT + 1)
		WHERE (C_W_ID, C_D_ID, C_ID) = (%d, %d, %d)
		RETURNING C_FIRST, C_MIDDLE, C_LAST, C_STREET_1, C_STREET_2, C_CITY, C_STATE, C_ZIP,
		C_PHONE, C_SINCE, C_CREDIT,C_CREDIT_LIM, C_DISCOUNT, C_BALANCE
	`, amount, warehouseID, districtID, customerID)

	if err := tx.QueryRow(sqlStatement).Scan(
		&customer.FirstName,
		&customer.MiddleName,
		&customer.LastName,
		&customer.Addr.Street1,
		&customer.Addr.Street2,
		&customer.Addr.City,
		&customer.Addr.State,
		&customer.Addr.Zip,
		&customer.Phone,
		&customer.Since,
		&customer.Credit,
		&customer.CreditLimit,
		&customer.Discount,
		&customer.Balance); err != nil {
		return nil, fmt.Errorf("error occurred in updating the customer details. Err: %v", err)
	}

	return
}
