package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	nodeID         = flag.Int("node", 0, "Pass the node id between 1 and 5 to start")
	env            = flag.String("env", "", "Pass the environment type to run: \"dev\" or \"prod\"")
	configFilePath = flag.String("config", "", "Pass the configuration file that contains the host and peer addresses")
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

func init() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		panic("use the flags before the command ")
	}
	// if flag.Args()[0] != "start" {
	// 	panic("pass a command: \"start\" to start the cluster")
	// }

	if *env == "prod" {
		if *nodeID < 1 || *nodeID > 5 {
			panic("pass a correct node id [1, 5] via the -node flag")
		}
		if len(*configFilePath) == 0 {
			panic("pass a correct config file path via the -config flag")
		}
	} else {
		*env = "dev"
		*nodeID = 1
		*configFilePath = "configs/dev/setup.json"
		fmt.Println("Using \"dev\" settings with node id = 1")
	}
}

func main() {
	function := flag.Args()[0]
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
	var cmd exec.Cmd

	switch function {
	case "start":
		cmd = start(config)
	case "stop":
		cmd = execute(config, function)
	case "sql":
		cmd = execute(config, function)
	case "init":
		cmd = execute(config, function)
	case "load":
		load(config)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)
}

func start(config configuration) exec.Cmd {
	var hostNode node
	joinNodes := make([]string, len(config.Nodes))

	for key, value := range config.Nodes {
		if value.ID == *nodeID {
			hostNode = value
		}
		joinNodes[key] = fmt.Sprintf("%s:%d", value.Host, value.Port)
	}

	cmd := &exec.Cmd{
		Path: "scripts/server.sh",
		Args: []string{"scripts/server.sh",
			"start",
			fmt.Sprintf("%s/cdb-server/node-files/node%d", config.WorkingDir, *nodeID),
			fmt.Sprintf("node%d", *nodeID),
			fmt.Sprintf("%s:%d", hostNode.Host, hostNode.Port),
			hostNode.HTTPAddr,
			strings.Join(joinNodes, ","),
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Dir:    ".",
	}

	return *cmd
}

func execute(config configuration, funcName string) exec.Cmd {
	var hostNode node

	for _, value := range config.Nodes {
		if value.ID == *nodeID {
			hostNode = value
			break
		}
	}

	cmd := &exec.Cmd{
		Path: "scripts/server.sh",
		Args: []string{"scripts/server.sh",
			funcName,
			fmt.Sprintf("%s:%d", hostNode.Host, hostNode.Port),
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Dir:    ".",
	}

	return *cmd
}

func load(config configuration) {
	var hostNode node

	for _, value := range config.Nodes {
		if value.ID == *nodeID {
			hostNode = value
			break
		}
	}
	
}
