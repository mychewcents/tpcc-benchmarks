package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var (
	configFilePath = flag.String("config", "", "Configuration file path for the server")
	nodeID         = flag.Int("node", 0, "Node ID to be used to connect to")
	env            = flag.String("env", "dev", "Provide an env: \"dev\" or \"prod\"")
)

type configuration struct {
	DownloadURL string `json:"data_files_url"`
	WorkingDir  string `json:"working_dir"`
	Nodes       []node `json:"nodes"`
}

type node struct {
	ID       int    `json:"id"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	HTTPAddr string `json:"http_addr"`
}

func main() {
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

	configFile, err := os.Open(*configFilePath)
	if err != nil {
		panic("file cannot be read")
	}

	byteValue, _ := ioutil.ReadAll(configFile)

	var config configuration

	if err = json.Unmarshal(byteValue, &config); err != nil {
		panic(err)
	}

	configFile.Close()

	cmd := &exec.Cmd{
		Path:   "scripts/init_setup.sh",
		Args:   []string{"scripts/init_setup.sh", *env, config.WorkingDir, config.DownloadURL, fmt.Sprintf("node%d", nodeID)},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Dir:    ".",
	}

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)
}
