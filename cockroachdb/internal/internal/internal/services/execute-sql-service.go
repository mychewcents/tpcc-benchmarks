package services

import "database/sql"

type ExecuteSQLService interface {
	Execute(sqlFilePath string) bool
	ExecutePartitions(sqlFilePath string) bool
}

type executeSQLServiceImpl struct {
	db *sql.DB
}

func CreateExecuteSQLService(db *sql.DB) ExecuteSQLService {
	return &executeSQLServiceImpl{db: db}
}

func (eqs *executeSQLServiceImpl) Execute(sqlFilePath string) bool {
	return true
}

func (eqs *executeSQLServiceImpl) ExecutePartitions(sqlFilePath string) bool {

	return true
}
