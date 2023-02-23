package scanners

import (
	"database/sql"

	"stories/database"
	"stories/types"
)

// String a special string type which understands jsons `null'
type String struct {
	Value types.String
}

// to help scanning a single string in some cases
// ScanRow implements the database.Scanner interface which scans database rows into custom types
func (scanner *String) ScanRow(row database.RowInterface) (interface{}, error) {
	var value sql.NullString
	var err error
	err = row.Scan(&value)
	if err != nil {
		return nil, err
	}
	if value.Valid {
		scanner.Value = types.NewString(value.String)
		return scanner.Value, nil
	}
	scanner.Value = types.NilString()
	return scanner.Value, nil
}
