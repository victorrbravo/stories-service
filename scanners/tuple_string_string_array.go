package scanners

import (
	"database/sql"
	"stories/database"
	"stories/types"
)

// TupleStringStringArray scans a row with two columns:
//   - a single string (e.g. ta_tracking.fk_status)
//   - a string array (e.g. ta_tracking.tags)
type TupleStringStringArray struct {
	First  types.String
	Second types.StringArray
}

func (scanner *TupleStringStringArray) ScanRow(row database.RowInterface) (interface{}, error) {
	var value0 sql.NullString
	var value1 types.StringArray
	err := row.Scan(&value0, &value1)
	if err != nil {
		return nil, err
	}
	scanner.First = types.NilString()
	if value0.Valid {
		scanner.First = types.NewString(value0.String)
	}
	scanner.Second = value1
	return *scanner, nil
}
