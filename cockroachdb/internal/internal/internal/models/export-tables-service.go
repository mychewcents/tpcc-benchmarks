package models

// TableExportParams stores the attributes required to export the tables
type TableExportParams struct {
	BaseSQLStatement string
	ExportPath       string
}
