package helper

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// CreateCmdObj create the command object to be executed
func CreateCmdObj(cliArgs []string) *exec.Cmd {
	cmd := &exec.Cmd{
		Path:   cliArgs[0],
		Args:   cliArgs,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Dir:    ".",
	}

	return cmd
}

// ExecuteCmd executes a command passed and returns when completed
func ExecuteCmd(cmd *exec.Cmd) error {
	log.Printf("Waiting for command to finish...")

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error occurred while starting the command. Err: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error occurred while waiting for the command to finish. Err: %v", err)
	}

	log.Printf("Successfully executed the command.")
	return nil
}
