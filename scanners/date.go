package scanners

import (
	"time"

	"stories/database"
)

type Date time.Time

// ScanRow implements the database.Scanner interface which scans database rows into custom types
func (date *Date) ScanRow(row database.RowInterface) (interface{}, error) {
	var err error
	err = row.Scan(date)
	return *date, err
}
