package main

import (
	"flag"

	caller "github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/common/logging"
)

var (
	functionName   = ""
	configFilePath = flag.String("config", "", "Configuration file path for the server")
	nodeID         = flag.Int("node", 0, "Node ID to be used to connect to")
	env            = flag.String("env", "", "Provide an env: \"dev\" or \"prod\"")
	experiment     = flag.Int("exp", 0, "Experiment Number")
)

func init() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		panic("use the flags before the command ")
	}

	if len(*configFilePath) == 0 {
		panic("provide a custom configuration file via -config flag")
	}
	if *env != "prod" && *env != "dev" {
		panic("provide the right environment via -env flag")
	}
	if *nodeID < 1 || *nodeID > 5 {
		panic("provide the right node id via -node flag")
	}
	functionName = flag.Args()[0]

	if functionName == "run-exp" {
		if *experiment == 0 {
			panic("provide Experiment and Client number to proceed")
		}
	}
	if err := logging.SetupLogOutput("server", "logs"); err != nil {
		panic(err)
	}
}

func main() {
	caller.ProcessServerSetupRequest(functionName, *configFilePath, *env, *nodeID, *experiment)
}
