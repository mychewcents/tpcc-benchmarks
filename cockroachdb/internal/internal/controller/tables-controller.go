package controller

import (
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/common/cdbconn"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/common/config"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/handler"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/internal/services"
)

type tablesControllerImpl struct {
	c config.Configuration
	e services.ExportTablesService
	p services.LoadProcessedTablesService
	r services.LoadRawTablesService
}

// CreateTablesController creates a new controller to export tables
func CreateTablesController(c config.Configuration) handler.NewTablesController {
	db, err := cdbconn.CreateConnection(c.HostNode)
	if err != nil {
		panic(err)
	}

	return &tablesControllerImpl{
		c: c,
		e: services.CreateExportTablesService(c, db),
		p: services.CreateLoadProcessedTablesService(db),
		r: services.CreateLoadRawTablesService(db),
	}
}

func (etc *tablesControllerImpl) ExportTables() (err error) {
	err = etc.e.Export()

	return
}

func (etc *tablesControllerImpl) LoadProcessedTables() (err error) {
	err = etc.p.Load()

	return
}

func (etc *tablesControllerImpl) LoadRawTables() (err error) {
	err = etc.r.Load()

	return
}
