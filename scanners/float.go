package scanners

import (
	"database/sql"

	"github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/types"
)

type Float struct {
	Value types.Float
}

// ScanRow implements the database.Scanner interface which scans database rows into custom types
func (scanner *Float) ScanRow(row database.RowInterface) (interface{}, error) {
	var value sql.NullFloat64
	var err error
	err = row.Scan(&value)
	if err != nil {
		return nil, err
	}
	if value.Valid {
		scanner.Value = types.NewFloat(value.Float64)
		return scanner.Value, nil
	}
	scanner.Value = types.NilFloat()
	return scanner, nil
}
