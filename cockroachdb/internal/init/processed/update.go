package processedtables

import (
	"log"

	"github.com/mychewcents/ddbms-project/cockroachdb/internal/connection/config"
)

// UpdateParent udpates parent tables
func UpdateParent(c config.Configuration) error {
	log.Printf("No update required")
	return nil
}

// UpdatePartitions updates partitions of the tables
func UpdatePartitions(c config.Configuration) error {
	log.Printf("No update required")
	return nil
}
