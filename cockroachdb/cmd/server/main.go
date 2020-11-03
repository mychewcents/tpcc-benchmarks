package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/cdbconn"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/init/tables"
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/logging"
)

var (
	functionName   = ""
	configFilePath = flag.String("config", "", "Configuration file path for the server")
	nodeID         = flag.Int("node", 0, "Node ID to be used to connect to")
	env            = flag.String("env", "", "Provide an env: \"dev\" or \"prod\"")
	experiment     = flag.Int("exp", 0, "Experiment Number")
	client         = flag.Int("client", 0, "Client Number")
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
		if *experiment == 0 || *client == 0 {
			panic("provide Experiment and Client number to proceed")
		}
	}
	if err := logging.SetupLogOutput("server", "logs"); err != nil {
		panic(err)
	}
}

func main() {
	c := config.GetConfig(*configFilePath, *nodeID)

	var cmd exec.Cmd

	switch functionName {
	case "start":
		cmd = start(c)
	case "stop":
		cmd = execute(c)
	case "init":
		cmd = execute(c)
	case "load":
		load(c)
	case "run-exp":
		cmd = run(c)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Err: %v", err)
	}
	log.Printf("Waiting for command to finish...")
	err := cmd.Wait()
	log.Printf("Command finished with error: %v", err)
}

func start(c config.Configuration) exec.Cmd {
	joinNodes := make([]string, len(c.Nodes))

	for key, value := range c.Nodes {
		joinNodes[key] = fmt.Sprintf("%s:%d", value.Host, value.Port)
	}

	cmd := &exec.Cmd{
		Path: "scripts/server.sh",
		Args: []string{"scripts/server.sh",
			"start",
			fmt.Sprintf("%s/cdb-server/node-files/%s", c.WorkingDir, c.HostNode.Name),
			fmt.Sprintf("%s:%d", c.HostNode.Host, c.HostNode.Port),
			c.HostNode.HTTPAddr,
			strings.Join(joinNodes, ","),
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Dir:    ".",
	}

	return *cmd
}

func execute(c config.Configuration) exec.Cmd {
	cmd := &exec.Cmd{
		Path: "scripts/server.sh",
		Args: []string{"scripts/server.sh",
			functionName,
			fmt.Sprintf("%s:%d", c.HostNode.Host, c.HostNode.Port),
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Dir:    ".",
	}

	return *cmd
}

func load(c config.Configuration) {

	db, err := cdbconn.CreateConnection(c.HostNode)
	if err != nil {
		panic("load function couldn't create a connection to the server")
	}

	fmt.Printf("Executing the SQL: scripts/sql/drop-partitions.sql")
	if err := tables.ExecuteSQLForPartitions(db, 10, 10, "scripts/sql/drop-partitions.sql"); err != nil {
		fmt.Println(err)
		return
	}

	sqlScripts := []string{
		"scripts/sql/drop-raw.sql",
		"scripts/sql/create-raw.sql",
		"scripts/sql/load-raw.sql",
		"scripts/sql/update-raw.sql",
	}

	for _, value := range sqlScripts {
		fmt.Printf("\nExecuting the SQL: %s", value)
		if err := tables.ExecuteSQL(db, value); err != nil {
			fmt.Println(err)
			return
		}
	}

	sqlScripts = []string{
		"scripts/sql/create-partitions.sql",
		"scripts/sql/load-partitions.sql",
		"scripts/sql/update-partitions.sql",
	}

	for _, value := range sqlScripts {
		fmt.Printf("\nExecuting the SQL: %s", value)
		if err := tables.ExecuteSQLForPartitions(db, 10, 10, value); err != nil {
			fmt.Println(err)
			return
		}
	}

	log.Println("Initialization Complete!")
	fmt.Println("\nInitialization Complete!")
}

func run(c config.Configuration) exec.Cmd {
	cmd := &exec.Cmd{
		Path: "scripts/run.sh",
		Args: []string{"scripts/run.sh",
			*env,
			fmt.Sprintf("%d", *experiment),
			fmt.Sprintf("%d", *nodeID),
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Dir:    ".",
	}

	return *cmd
}
