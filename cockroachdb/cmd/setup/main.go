package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type configuration struct {
	DownloadURL string `json:"data_files_url"`
	WorkingDir  string `json:"working_dir"`
	Nodes       []node `json:"nodes"`
}

type node struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func main() {
	var setupConfigFile string

	if len(os.Args) != 3 {
		fmt.Println("No custom configuration file. Using the \"dev\" configuration file: configs/dev/setup.json file.")
		setupConfigFile = "configs/dev/setup.json"
	} else if os.Args[1] == "dev" || os.Args[1] == "prod" {
		setupConfigFile = os.Args[2]
	} else {
		panic("use \"dev\" or \"prod\" as the first argument and pass the \"configuration file\" in the second arugment")
	}

	configFile, err := os.Open(setupConfigFile)
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
		Args:   []string{"scripts/init_setup.sh", os.Args[1], config.WorkingDir, config.DownloadURL},
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
