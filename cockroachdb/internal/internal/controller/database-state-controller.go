package controller

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/helper"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/services"
)

// DatabaseStateController interface to get the database state
type DatabaseStateController interface {
	CalculateDBState(experiment int, dirPath string) error
}

type databaseStateControllerImpl struct {
	s services.DatabaseStateService
}

// CreateDatabaseStateController creates a new controller to get the database state
func CreateDatabaseStateController(db *sql.DB) DatabaseStateController {
	return &databaseStateControllerImpl{
		s: services.CreateDatabaseStateService(db),
	}
}

func (dbsc *databaseStateControllerImpl) CalculateDBState(experiment int, dirPath string) (err error) {

	result, err := dbsc.s.CalculateDBState()
	if err != nil {
		return err
	}

	outputCSVString := fmt.Sprintf("%d,%f,%f,%d,%f,%f,%d,%d,%d,%d,%f,%d,%d,%f,%d,%d",
		experiment,
		result.TotalYTDWarehouse,
		result.TotalYTDDistrict,
		result.SumOrderIDs,
		result.CBalance,
		result.CYTDPayment,
		result.CPaymentCount,
		result.CDeliveryCount,
		result.MaxOrderID,
		result.TotalOrderLineCount,
		result.TotalOrderAmount,
		result.TotalQuantity,
		result.TotalStock,
		result.TotalYTDStock,
		result.TotalOrderCount,
		result.TotalRemoteOrderCount,
	)

	finalCSVPath := fmt.Sprintf("%s/%d.csv", dirPath, experiment)
	err = helper.WriteCSVFile(outputCSVString, finalCSVPath)
	if err != nil {
		log.Printf("DB State: %s", outputCSVString)
		log.Fatalf("error occurred in writing the csv file. Err: %v", err)
		return err
	}

	return nil
}
