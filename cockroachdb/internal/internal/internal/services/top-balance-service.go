package services

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dao"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
)

// TopBalanceService creates the top balance interface
type TopBalanceService interface {
	ProcessTransaction(req *models.TopBalance) (*models.TopBalanceOutput, error)
	Print(result *models.TopBalanceOutput)
}

type topBalanceServiceImpl struct {
	db *sql.DB
	c  dao.CustomerDao
	d  dao.DistrictDao
	w  dao.WarehouseDao
}

// CreateTopBalanceService creates the new service
func CreateTopBalanceService(db *sql.DB) TopBalanceService {
	return &topBalanceServiceImpl{
		db: db,
		c:  dao.CreateCustomerDao(db),
		d:  dao.CreateDistrictDao(db),
		w:  dao.CreateWarehouseDao(db),
	}
}

func (tbs *topBalanceServiceImpl) ProcessTransaction(req *models.TopBalance) (result *models.TopBalanceOutput, err error) {
	result, err = tbs.execute()
	if err != nil {
		return nil, err
	}

	return
}

func (tbs *topBalanceServiceImpl) execute() (result *models.TopBalanceOutput, err error) {
	customers, err := tbs.c.GetCustomersWithTopBalance(10)
	if err != nil {
		return nil, err
	}

	warehouseIDs := make([]int, 10)
	districtIDs := make([]int, 10)

	idx := 0
	for _, value := range customers {
		warehouseIDs[idx] = value.WarehouseID
		districtIDs[idx] = value.DistrictID
		idx++
	}

	districtNames, err := tbs.d.GetDistrictNames(warehouseIDs, districtIDs)
	if err != nil {
		return nil, err
	}

	warehouseNames, err := tbs.w.GetWarehouseNames(warehouseIDs)
	if err != nil {
		return nil, err
	}

	result.Rows = make([]*models.TopBalanceOutputRow, 10)

	idx = 0
	for _, value := range customers {
		result.Rows[idx] = &models.TopBalanceOutputRow{
			WarehouseID:   value.WarehouseID,
			DistrictID:    value.DistrictID,
			CustomerID:    value.CustomerID,
			FirstName:     value.FirstName,
			MiddleName:    value.MiddleName,
			LastName:      value.LastName,
			Balance:       value.Balance,
			WarehouseName: warehouseNames[value.WarehouseID],
			DistrictName:  districtNames[value.WarehouseID][value.DistrictID],
		}
		idx++
	}

	return
}

func (tbs *topBalanceServiceImpl) Print(result *models.TopBalanceOutput) {
	var topBalanceOutputBuilder strings.Builder

	for _, row := range result.Rows {
		topBalanceOutputBuilder.WriteString(fmt.Sprintf("Customer: %s %s %s, Balance: %f, Warehouse: %s, District: %s",
			row.FirstName, row.MiddleName, row.LastName, row.Balance, row.WarehouseName, row.DistrictName))
	}
}
