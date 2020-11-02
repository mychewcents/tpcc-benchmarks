package dbstate

import "database/sql"

// RecordDBState records the DB State for the experiment performed
func RecordDBState(db *sql.DB, experiment int, path string) {

}

func getWarehouseState(db *sql.DB) (float64, error) {
	var sumYTD float64

	sqlStatement := `SELECT sum(W_YTD) FROM District`
	row := db.QueryRow(sqlStatement)
	if err := row.Scan(&sumYTD); err != nil {
		return 0.0, err
	}
	return sumYTD, nil
}
