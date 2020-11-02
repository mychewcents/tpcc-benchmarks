package logging

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

// SetupLogOutput sets the log output file
func SetupLogOutput(logType string, path string, fileArgs ...int) error {
	var fileName string
	if logType == "init" {
		fileName = fmt.Sprintf("%s/init_%s.log", path, time.Now())
	} else if logType == "run" {
		fileName = fmt.Sprintf("%s/exp_%d_client_%d_%s.log", path, fileArgs[0], fileArgs[1], time.Now())
	} else {
		return errors.New("no matching log type passed")
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}

	log.SetOutput(file)
	return nil
}
