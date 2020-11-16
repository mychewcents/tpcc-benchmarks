package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/common/config"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/helper"
)

// Start starts the server on the host machine
func Start(c config.Configuration) {
	joinNodes := make([]string, len(c.Nodes))

	for key, value := range c.Nodes {
		joinNodes[key] = fmt.Sprintf("%s:%d", value.Host, value.Port)
	}

	cliArgs := []string{"scripts/cli/server.sh",
		"start",
		fmt.Sprintf("%s/cdb-server/node-files/%s", c.WorkingDir, c.HostNode.Name),
		fmt.Sprintf("%s:%d", c.HostNode.Host, c.HostNode.Port),
		c.HostNode.HTTPAddr,
		strings.Join(joinNodes, ","),
	}

	execute(cliArgs)
}

// Initialize executes the init call to CRDB
func Initialize(c config.Configuration) {
	cliArgs := []string{"scripts/cli/server.sh",
		"init",
		fmt.Sprintf("%s:%d", c.HostNode.Host, c.HostNode.Port),
	}

	execute(cliArgs)
}

// Stop stops the server on the host machine
func Stop(c config.Configuration) {
	cliArgs := []string{"scripts/cli/server.sh",
		"stop",
		fmt.Sprintf("%s:%d", c.HostNode.Host, c.HostNode.Port),
	}

	execute(cliArgs)
}

// RunExperiments runs the TPCC experiment as per the project
func RunExperiments(configFilePath, env string, nodeID, experiment int) {
	cliArgs := []string{"scripts/cli/run.sh",
		env,
		fmt.Sprintf("%d", experiment),
		fmt.Sprintf("%d", nodeID),
		configFilePath,
	}

	execute(cliArgs)
}

func execute(cliArgs []string) {
	if err := helper.ExecuteCmd(helper.CreateCmdObj(cliArgs)); err != nil {
		log.Fatalf("error occurred. Err: %v", err)
	}
}
