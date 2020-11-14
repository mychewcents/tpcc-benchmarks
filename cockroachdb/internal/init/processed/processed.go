package processedtables

import "github.com/mychewcents/tpcc-benchmarks/cockroachdb/internal/common/config"

// PerformETL performs ETL for processed tables
func PerformETL(c config.Configuration) error {
	if err := DropPartitions(c); err != nil {
		return err
	}
	if err := DropParent(c); err != nil {
		return err
	}

	if err := CreateParent(c); err != nil {
		return err
	}
	if err := CreatePartitions(c); err != nil {
		return err
	}

	if err := LoadParent(c); err != nil {
		return err
	}
	if err := LoadPartitions(c); err != nil {
		return err
	}

	if err := UpdateParent(c); err != nil {
		return err
	}
	if err := UpdatePartitions(c); err != nil {
		return err
	}

	return nil
}
