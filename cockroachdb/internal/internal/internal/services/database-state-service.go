package services

import (
	"database/sql"
	"fmt"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dao"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
)

// DatabaseStateService interface to call the service to calculate the DB State
type DatabaseStateService interface {
	CalculateDBState() (*models.DatabaseState, error)
}

type databaseStateServiceImpl struct {
	db *sql.DB
	d  dao.DistrictDao
	c  dao.CustomerDao
	o  dao.OrderDao
	ol dao.OrderLineDao
	s  dao.StockDao
}

// CreateDatabaseStateService creates the service to calculate the database state
func CreateDatabaseStateService(db *sql.DB) DatabaseStateService {
	return &databaseStateServiceImpl{
		db: db,
		d:  dao.CreateDistrictDao(db),
		c:  dao.CreateCustomerDao(db),
		o:  dao.CreateOrderDao(db),
		ol: dao.CreateOrderLineDao(db),
		s:  dao.CreateStockDao(db),
	}
}

func (dbss *databaseStateServiceImpl) CalculateDBState() (result *models.DatabaseState, err error) {
	result.TotalYTDDistrict, result.SumOrderIDs, err = dbss.d.GetFinalState()
	if err != nil {
		return nil, fmt.Errorf("error occurred while recording the district and warehouse state. Err: %v", err)
	}

	result.TotalYTDWarehouse = result.TotalYTDDistrict

	result.CBalance, result.CYTDPayment, result.CPaymentCount, result.CDeliveryCount, err = dbss.c.GetFinalState()
	if err != nil {
		return nil, fmt.Errorf("error occurred while recording the customer state. Err: %v", err)
	}

	result.MaxOrderID, result.TotalOrderLineCount, result.TotalOrderAmount, err = dbss.o.GetFinalState()
	if err != nil {
		return nil, fmt.Errorf("error occurred while recording the order state. Err: %v", err)
	}

	result.TotalQuantity, err = dbss.ol.GetFinalState()
	if err != nil {
		return nil, fmt.Errorf("error occurred while recording the order line state. Err: %v", err)
	}

	result.TotalStock, result.TotalOrderCount, result.TotalRemoteOrderCount, result.TotalYTDStock, err = dbss.s.GetFinalState()
	if err != nil {
		return nil, fmt.Errorf("error occurred while recording the stock state. Err: %v", err)
	}

	return
}
