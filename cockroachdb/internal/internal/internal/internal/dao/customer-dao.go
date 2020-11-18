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
	GetCustomersWithTopBalance(num int) ([]*dbdatamodel.Customer, error)
	DeliverOrder(tx *sql.Tx, warehouseID, districtID, customerID int, totalAmount float64) error
	GetFinalState() (float64, float64, int, int, error)
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
	sqlStatement := fmt.Sprintf("SELECT C_FIRST, C_MIDDLE, C_LAST, C_CREDIT, C_DISCOUNT, C_BALANCE FROM CUSTOMER WHERE C_W_ID = $1 AND C_D_ID = $2 AND C_ID = $3")

	var firstName, middleName, lastName, credit string
	var discount, balance float64

	row := cs.db.QueryRow(sqlStatement, warehouseID, districtID, customerID)
	if err := row.Scan(&firstName, &middleName, &lastName, &credit, &discount, &balance); err != nil {
		return nil, fmt.Errorf("error occured in getting the customer details. Err: %v", err)
	}

	customer = &dbdatamodel.Customer{
		WarehouseID: warehouseID,
		DistrictID:  districtID,
		CustomerID:  customerID,
		FirstName:   firstName,
		MiddleName:  middleName,
		LastName:    lastName,
		Credit:      credit,
		Discount:    discount,
		Balance:     balance,
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

// GetCustomerWithTopBalance returns the top number of balance Customers
func (cs *customerDaoImpl) GetCustomersWithTopBalance(num int) (result []*dbdatamodel.Customer, err error) {
	result = make([]*dbdatamodel.Customer, num)

	sqlStatement := fmt.Sprintf("SELECT C_FIRST, C_MIDDLE, C_LAST, C_W_ID, C_D_ID, C_BALANCE FROM CUSTOMER ORDER BY C_BALANCE DESC LIMIT %d", num)

	rows, err := cs.db.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading the customer details. Err: %v", err)
	}
	defer rows.Close()

	var firstName, middleName, lastName string
	var wID, dID, cID int
	var balance float64

	idx := 0
	for rows.Next() {
		if err := rows.Scan(&firstName, &middleName, &lastName, &wID, &dID, &balance); err != nil {
			return nil, fmt.Errorf("error occured while scanning the customer details. Err: %v", err)
		}

		result[idx] = &dbdatamodel.Customer{
			WarehouseID: wID,
			DistrictID:  dID,
			CustomerID:  cID,
			FirstName:   firstName,
			MiddleName:  middleName,
			LastName:    lastName,
			Balance:     balance,
		}
		idx++
	}

	return
}

// DeliverOrder delivers the order and updates the customer with balance
func (cs *customerDaoImpl) DeliverOrder(tx *sql.Tx, warehouseID, districtID, customerID int, amount float64) error {
	sqlStatement := fmt.Sprintf("UPDATE CUSTOMER SET (C_BALANCE, C_DELIVERY_CNT) = (C_BALANCE + %f, C_DELIVERY_CNT + 1) WHERE C_W_ID=%d AND C_D_ID=%d AND C_ID=%d", amount, warehouseID, districtID, customerID)

	if _, err := tx.Exec(sqlStatement); err != nil {
		return err
	}

	return nil
}

func (cs *customerDaoImpl) GetFinalState() (balance, ytdPayment float64, paymentCount, deliveryCount int, err error) {
	sqlStatement := "SELECT SUM(C_BALANCE), SUM(C_YTD_PAYMENT), SUM(C_PAYMENT_CNT), SUM(C_DELIVERY_CNT) FROM Customer"

	row := cs.db.QueryRow(sqlStatement)
	if err := row.Scan(&balance, &ytdPayment, &paymentCount, &deliveryCount); err != nil {
		return 0.0, 0.0, 0, 0, err
	}

	return
}
