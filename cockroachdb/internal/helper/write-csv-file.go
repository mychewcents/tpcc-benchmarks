package helper

import (
	"fmt"
	"os"
)

// WriteCSVFile writes the CSV file
func WriteCSVFile(csvString string, fileName string) error {
	csvFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("Error in creating CSV file, Err: %v", err)
	}
	defer csvFile.Close()

	if _, err := csvFile.WriteString(csvString); err != nil {
		return fmt.Errorf("Error in writing the db state CSV file, Err: %v", err)
	}

	return nil
}
