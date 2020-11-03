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
		fileName = fmt.Sprintf("%s/init_%d.log", path, time.Now().Unix())
	} else if logType == "exp" {
		fileName = fmt.Sprintf("%s/exp_%d_client_%d_%s.log", path, fileArgs[0], fileArgs[1], time.Now())
	} else if logType == "dbstate" {
		fileName = fmt.Sprintf("%s/dbstate_%d_%s.log", path, fileArgs[0], time.Now())
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
