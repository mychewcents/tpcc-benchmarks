package router

import (
	"log"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/common/cdbconn"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/common/config"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/controller"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/handler"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/server"
)

// ServerSetupRouter stores the handlers
type ServerSetupRouter struct {
	c        config.Configuration
	handlers map[string]handler.NewLoadTablesController
}

// CreateServerSetupRouter creates a new router to process server requests
func CreateServerSetupRouter(configFilePath string, nodeID int) *ServerSetupRouter {
	c := config.GetConfig(configFilePath, nodeID)

	handlers := make(map[string]handler.NewLoadTablesController)
	db, err := cdbconn.CreateConnection(c.HostNode)
	if err != nil {
		panic(err)
	}

	handlers["LoadProcessedTables"] = controller.CreateLoadProcessedTablesController(db)
	handlers["LoadRawTables"] = controller.CreateLoadRawTablesController(db)

	return &ServerSetupRouter{c: c, handlers: handlers}
}

// ProcessServerSetupRequest processes each transaction upon input
func (ssr *ServerSetupRouter) ProcessServerSetupRequest(functionName string, configFilePath, env string, nodeID, experiment int) {
	switch functionName {
	case "start":
		server.Start(ssr.c)
	case "stop":
		server.Stop(ssr.c)
	case "init":
		server.Initialize(ssr.c)
	case "load":
		if err := ssr.handlers["LoadRawTables"].LoadTables(); err != nil {
			log.Fatalf("Err: %v", err)
		}
	case "load-csv":
		if err := ssr.handlers["LoadProcessedTables"].LoadTables(); err != nil {
			log.Fatalf("Err: %v", err)
		}
	case "run-exp":
		server.RunExperiments(configFilePath, env, nodeID, experiment)
	case "export":
		// exportCSV(c)
	}
}
