package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Configuration creates the configuration object
type Configuration struct {
	DownloadURL string   `json:"data_files_url"`
	WorkingDir  string   `json:"working_dir"`
	Nodes       []Server `json:"nodes"`
	HostNode    Server
}

// Server denotes the server configurations to connect to
type Server struct {
	ID       int    `json:"id"`
	Host     string `json:"host"`
	Name     string
	Port     int    `json:"port"`
	HTTPAddr string `json:"http_addr"`
	Database string `json:"database"`
	Username string `json:"username"`
}

// GetConfig reads the input configuration file and returns the object
func GetConfig(path string, nodeID int) Configuration {

	configFile, err := os.Open(path)
	if err != nil {
		panic("file cannot be read")
	}
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)

	var c Configuration

	if err = json.Unmarshal(byteValue, &c); err != nil {
		panic(err)
	}

	for _, value := range c.Nodes {
		if value.ID == nodeID {
			c.HostNode = value
			c.HostNode.Name = fmt.Sprintf("node%d", value.ID)
			break
		}
		value.Name = fmt.Sprintf("node%d", value.ID)
	}

	return c
}
