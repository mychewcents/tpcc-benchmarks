package services

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/internal/dao"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/models"
)

// StockLevelService provides the service to StockLevel transaction
type StockLevelService interface {
	ProcessTransaction(req *models.StockLevel) (*models.StockLevelOutput, error)
}

type stockLevelServiceImpl struct {
	db *sql.DB
	d  dao.DistrictDao
	s  dao.StockDao
}

// CreateStockLevelService creates the new stock level service
func CreateStockLevelService(db *sql.DB) StockLevelService {
	return &stockLevelServiceImpl{
		db: db,
		d:  dao.CreateDistrictDao(db),
		s:  dao.CreateStockDao(db),
	}
}

// ProcessTransaction processes the stock level transaction
func (sls *stockLevelServiceImpl) ProcessTransaction(req *models.StockLevel) (*models.StockLevelOutput, error) {
	log.Printf("Starting the Stock Level Transaction for: w=%d d=%d t=%d n=%d", req.WarehouseID, req.DistrictID, req.Threshold, req.LastNOrders)
	result, err := sls.execute(req)
	if err != nil {
		return nil, fmt.Errorf("error occured while executing the stock level transaction. Err: %v", err)
	}

	log.Printf("Completed the Stock Level Transaction for: w=%d d=%d t=%d n=%d", req.WarehouseID, req.DistrictID, req.Threshold, req.LastNOrders)
	return result, nil
}

func (sls *stockLevelServiceImpl) execute(req *models.StockLevel) (result *models.StockLevelOutput, err error) {
	result.LastOrderID, err = sls.d.GetLastOrderID(req.WarehouseID, req.DistrictID)
	if err != nil {
		return nil, err
	}

	result.StartOrderID = result.LastOrderID - req.LastNOrders

	result.TotalItems, err = sls.s.GetStockItemsBelowThreshold(req.WarehouseID, req.DistrictID, req.Threshold, result.StartOrderID, result.LastOrderID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (sls *stockLevelServiceImpl) Print(o *models.StockLevelOutput) {
	fmt.Println(fmt.Sprintf("Total Number of Items below threshold: %d , for Order IDs between %d - %d", o.TotalItems, o.StartOrderID, o.LastOrderID))
}
