package main

import (
	"flag"
	"log"
	"os"
	"os/exec"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
)

var (
	configFilePath = flag.String("config", "", "Configuration file path for the server")
	nodeID         = flag.Int("node", 0, "Node ID to be used to connect to")
	env            = flag.String("env", "dev", "Provide an env: \"dev\" or \"prod\"")
)

func init() {
	flag.Parse()

	if len(*configFilePath) == 0 {
		panic("provide a custom configuration file via -config flag")
	}
	if *env != "prod" && *env != "dev" {
		panic("provide the right environment via -env flag")
	}
	if *nodeID < 1 || *nodeID > 5 {
		panic("provide the right node id via -node flag")
	}
}

func main() {
	c := config.GetConfig(*configFilePath, *nodeID)

	cmd := &exec.Cmd{
		Path:   "scripts/init_setup.sh",
		Args:   []string{"scripts/init_setup.sh", *env, c.WorkingDir, c.DownloadURL, c.HostNode.Name},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Dir:    ".",
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Err: %v", err)
	}
	log.Printf("Waiting for command to finish...")
	if err := cmd.Wait(); err != nil {
		log.Fatalf("Err: %v", err)
	}
	log.Printf("Command finished")
}
