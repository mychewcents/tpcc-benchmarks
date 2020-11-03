package main

import (
	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
)

func export(c config.Configuration) {

	// cmd := &exec.Cmd{
	// 	Path: "scripts/export_data.sh",
	// 	Args: []string{"scripts/export_data.sh",
	// 		"start",
	// 		fmt.Sprintf("%s/cdb-server/node-files/%s", c.WorkingDir, c.HostNode.Name),
	// 		fmt.Sprintf("%s:%d", c.HostNode.Host, c.HostNode.Port),
	// 		c.HostNode.HTTPAddr,
	// 		strings.Join(joinNodes, ","),
	// 	},
	// 	Stdout: os.Stdout,
	// 	Stderr: os.Stderr,
	// 	Dir:    ".",
	// }
}
