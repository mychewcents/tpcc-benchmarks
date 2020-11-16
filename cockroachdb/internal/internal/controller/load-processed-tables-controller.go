package controller

import (
	"database/sql"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/services"
)

// LoadProcessedTablesController interface to the processed tables loading controller
type LoadProcessedTablesController interface {
	LoadTables(configFilePath string) bool
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

func (lptc *loadProcessedTablesControllerImpl) LoadTables(configFilePath string) bool {

	return true
}
