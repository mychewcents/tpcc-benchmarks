package helper

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// ReadFile reads the file and returns the contents as string
func ReadFile(filePath string) (fileContents string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error occurred while reading the file: %s. Err: %v", filePath, err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var fileContentsBuilder strings.Builder
	fileContentsBuilder.WriteString(string(byteValue))

	fileContents = fileContentsBuilder.String()

	return
}

// WriteFile writes the CSV file
func WriteFile(contentString string, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error in creating file. Err: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(contentString); err != nil {
		return fmt.Errorf("Error in writing the db state CSV file, Err: %v", err)
	}

	return nil
}
