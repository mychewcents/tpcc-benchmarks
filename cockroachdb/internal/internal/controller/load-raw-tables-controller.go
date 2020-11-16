package controller

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/helper"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/services"
)

// LoadRawTablesController interface to the raw tables controller
type LoadRawTablesController interface {
	LoadTables() error
}

type loadRawTablesControllerImpl struct {
	s  services.ExecuteSQLService
	ci services.CustomerItemsPairService
}

// CreateLoadRawTablesController creates new controller for the raw tables
func CreateLoadRawTablesController(db *sql.DB) LoadRawTablesController {
	return &loadRawTablesControllerImpl{
		s:  services.CreateExecuteSQLService(db),
		ci: services.CreateCustomerItemsPairService(db),
	}
}

func (lrtc *loadRawTablesControllerImpl) LoadTables() (err error) {
	var sqlString string
	log.Println("Dropping tables...")

	sqlString, err = helper.ReadFile("scripts/sql/raw/drop-partitions.sql")
	if err != nil {
		return err
	}
	if err := lrtc.s.ExecutePartitions(10, 10, sqlString); err != nil {
		return fmt.Errorf("error occured while dropping partitions. Err: %v", err)
	}

	sqlString, err = helper.ReadFile("scripts/sql/raw/drop.sql")
	if err != nil {
		return err
	}
	if err := lrtc.s.Execute(sqlString); err != nil {
		return fmt.Errorf("error occured while dropping parent tables. Err: %v", err)
	}

	log.Println("Creating tables...")

	sqlString, err = helper.ReadFile("scripts/sql/raw/create.sql")
	if err != nil {
		return err
	}
	if err := lrtc.s.Execute(sqlString); err != nil {
		return fmt.Errorf("error occured while creating parent tables. Err: %v", err)
	}

	sqlString, err = helper.ReadFile("scripts/sql/raw/create-partitions.sql")
	if err != nil {
		return err
	}
	if err := lrtc.s.ExecutePartitions(10, 10, sqlString); err != nil {
		return fmt.Errorf("error occured while creating partitions. Err: %v", err)
	}

	log.Println("Loading tables...")

	sqlString, err = helper.ReadFile("scripts/sql/raw/load.sql")
	if err != nil {
		return err
	}
	if err := lrtc.s.Execute(sqlString); err != nil {
		return fmt.Errorf("error occured while loading parent tables. Err: %v", err)
	}

	sqlString, err = helper.ReadFile("scripts/sql/raw/load-partitions.sql")
	if err != nil {
		return err
	}
	if err := lrtc.s.ExecutePartitions(10, 10, sqlString); err != nil {
		return fmt.Errorf("error occured while loading partitions. Err: %v", err)
	}

	if err := lrtc.ci.LoadInitial(10, 10); err != nil {
		return fmt.Errorf("error occured while loading the customer items pair table. Err: %v", err)
	}

	log.Println("Updating tables...")

	sqlString, err = helper.ReadFile("scripts/sql/raw/update.sql")
	if err != nil {
		return err
	}
	if err := lrtc.s.Execute(sqlString); err != nil {
		return fmt.Errorf("error occured while updating parent tables. Err: %v", err)
	}

	sqlString, err = helper.ReadFile("scripts/sql/raw/update-partitions.sql")
	if err != nil {
		return err
	}
	if err := lrtc.s.ExecutePartitions(10, 10, sqlString); err != nil {
		return fmt.Errorf("error occured while updating partitions. Err: %v", err)
	}

	return nil
}
