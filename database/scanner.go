package database

// Scanner is used to scan full rows from a database query result
type Scanner interface {
	ScanRow(row RowInterface) (interface{}, error)
}
