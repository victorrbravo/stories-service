package scanners

import (
	"database/sql"

	"stories/database"
	"stories/types"
)

type Int struct {
	Value types.Int
}

// ScanRow implements the database.Scanner interface which scans database rows into custom types
func (scanner *Int) ScanRow(row database.RowInterface) (interface{}, error) {
	var value sql.NullInt64
	var err error
	err = row.Scan(&value)
	if err != nil {
		return nil, err
	}
	if value.Valid {
		scanner.Value = types.NewInt(int(value.Int64))
		return scanner.Value, nil
	}
	scanner.Value = types.NilInt
	return scanner, nil
}
