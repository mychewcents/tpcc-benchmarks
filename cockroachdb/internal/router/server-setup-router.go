package router

import (
	"log"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/common/config"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/controller"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/internal/handler"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/server"
)

// ServerSetupRouter stores the handlers
type ServerSetupRouter struct {
	c        config.Configuration
	handlers handler.NewTablesController
}

// CreateServerSetupRouter creates a new router to process server requests
func CreateServerSetupRouter(configFilePath string, nodeID int) *ServerSetupRouter {
	c := config.GetConfig(configFilePath, nodeID)

	handlers := controller.CreateTablesController(c)

	return &ServerSetupRouter{c: c, handlers: handlers}
}

// ProcessServerSetupRequest processes each transaction upon input
func (ssr *ServerSetupRouter) ProcessServerSetupRequest(functionName string, configFilePath, env string, nodeID, experiment int) {
	switch functionName {
	case "download-dataset":
		server.DownloadDataset(ssr.c)
	case "setup-dir":
		server.SetupDirectories(ssr.c, env)
	case "start":
		server.Start(ssr.c)
	case "stop":
		server.Stop(ssr.c)
	case "init":
		server.InitCluster(ssr.c)
	case "run-exp":
		server.RunExperiments(configFilePath, env, nodeID, experiment)
	case "load":
		if err := ssr.handlers.LoadRawTables(); err != nil {
			log.Fatalf("Err: %v", err)
		}
	case "load-csv":
		if err := ssr.handlers.LoadProcessedTables(); err != nil {
			log.Fatalf("Err: %v", err)
		}
	case "export":
		if err := ssr.handlers.ExportTables(); err != nil {
			log.Fatalf("Err: %v", err)
		}
	}
}
