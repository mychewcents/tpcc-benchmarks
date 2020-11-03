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
	if logType == "server" || logType == "setup" {
		fileName = fmt.Sprintf("%s/%s_%d.log", path, logType, time.Now().Unix())
	} else if logType == "exp" {
		fileName = fmt.Sprintf("%s/exp_%d_client_%d_%d.log", path, fileArgs[0], fileArgs[1], time.Now().Unix())
	} else if logType == "dbstate" {
		fileName = fmt.Sprintf("%s/dbstate_%d_%d.log", path, fileArgs[0], time.Now().Unix())
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
