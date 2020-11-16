package controller

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/helper"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/services"
)

// LoadProcessedTablesController interface to the processed tables loading controller
type LoadProcessedTablesController interface {
	LoadTables() error
}

type loadProcessedTablesControllerImpl struct {
	s services.ExecuteSQLService
}

// CreateLoadProcessedTablesController creates the load processed tables controller
func CreateLoadProcessedTablesController(db *sql.DB) LoadProcessedTablesController {
	return &loadProcessedTablesControllerImpl{
		s: services.CreateExecuteSQLService(db),
	}
}

func (lptc *loadProcessedTablesControllerImpl) LoadTables() (err error) {
	var sqlString string
	log.Println("Dropping tables...")

	sqlString, err = helper.ReadFile("scripts/sql/processed/drop-partitions.sql")
	if err != nil {
		return err
	}
	if err := lptc.s.ExecutePartitions(10, 10, sqlString); err != nil {
		return fmt.Errorf("error occured while dropping partitions. Err: %v", err)
	}

	sqlString, err = helper.ReadFile("scripts/sql/processed/drop.sql")
	if err != nil {
		return err
	}
	if err := lptc.s.Execute(sqlString); err != nil {
		return fmt.Errorf("error occured while dropping parent tables. Err: %v", err)
	}

	log.Println("Creating tables...")

	sqlString, err = helper.ReadFile("scripts/sql/processed/create.sql")
	if err != nil {
		return err
	}
	if err := lptc.s.Execute(sqlString); err != nil {
		return fmt.Errorf("error occured while creating parent tables. Err: %v", err)
	}

	sqlString, err = helper.ReadFile("scripts/sql/processed/create-partitions.sql")
	if err != nil {
		return err
	}
	if err := lptc.s.ExecutePartitions(10, 10, sqlString); err != nil {
		return fmt.Errorf("error occured while creating partitions. Err: %v", err)
	}

	log.Println("Loading tables...")

	sqlString, err = helper.ReadFile("scripts/sql/processed/load.sql")
	if err != nil {
		return err
	}
	if err := lptc.s.Execute(sqlString); err != nil {
		return fmt.Errorf("error occured while loading parent tables. Err: %v", err)
	}

	sqlString, err = helper.ReadFile("scripts/sql/processed/load-partitions.sql")
	if err != nil {
		return err
	}
	if err := lptc.s.ExecutePartitions(10, 10, sqlString); err != nil {
		return fmt.Errorf("error occured while loading partitions. Err: %v", err)
	}

	log.Println("Updating tables...")

	sqlString, err = helper.ReadFile("scripts/sql/processed/update-partitions.sql")
	if err != nil {
		return err
	}
	if err := lptc.s.ExecutePartitions(10, 10, sqlString); err != nil {
		return fmt.Errorf("error occured while updating partitions. Err: %v", err)
	}

	return nil
}
