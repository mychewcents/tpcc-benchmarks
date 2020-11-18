package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/common/config"
	"github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/helper"
)

// DownloadDataset downloads the project files as a zip and extracts it
func DownloadDataset(c config.Configuration) {
	cliArgs := []string{"scripts/cli/server.sh", "download-dataset", c.DownloadURL}

	err := Execute(cliArgs)
	log.Fatalf("error occurred. Err: %v", err)
}

// SetupDirectories sets up the directories required for further setup
func SetupDirectories(c config.Configuration, env string) {
	cliArgs := []string{"scripts/cli/server.sh", "setup-dirs", c.WorkingDir, c.HostNode.Name}

	err := Execute(cliArgs)
	log.Fatalf("error occurred. Err: %v", err)
}

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

	err := Execute(cliArgs)
	log.Fatalf("error occurred. Err: %v", err)
}

// InitCluster executes the init call to CRDB
func InitCluster(c config.Configuration) {
	cliArgs := []string{"scripts/cli/server.sh",
		"init",
		fmt.Sprintf("%s:%d", c.HostNode.Host, c.HostNode.Port),
	}

	err := Execute(cliArgs)
	log.Fatalf("error occurred. Err: %v", err)
}

// Stop stops the server on the host machine
func Stop(c config.Configuration) {
	cliArgs := []string{"scripts/cli/server.sh",
		"stop",
		fmt.Sprintf("%s:%d", c.HostNode.Host, c.HostNode.Port),
	}

	err := Execute(cliArgs)
	log.Fatalf("error occurred. Err: %v", err)
}

// RunExperiments runs the TPCC experiment as per the project
func RunExperiments(configFilePath, env string, nodeID, experiment int) {
	cliArgs := []string{"scripts/cli/run.sh",
		env,
		fmt.Sprintf("%d", experiment),
		fmt.Sprintf("%d", nodeID),
		configFilePath,
	}

	err := Execute(cliArgs)
	log.Fatalf("error occurred. Err: %v", err)
}

// Execute executes the shell command with passed cli arguments
func Execute(cliArgs []string) (err error) {
	err = helper.ExecuteCmd(helper.CreateCmdObj(cliArgs))

	return
}
