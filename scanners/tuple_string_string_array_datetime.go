package scanners

import (
	"database/sql"
	"github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/types"
)

// TupleStringStringArrayDateTime scans a row with three columns:
//   - a single string (e.g. ta_tracking.fk_status)
//   - a string array (e.g. ta_tracking.tags)
//   - a timestamp (e.g. ta_tracking.updated_at)
type TupleStringStringArrayDateTime struct {
	First  types.String
	Second types.StringArray
	Third  types.DateTime
}

func (scanner *TupleStringStringArrayDateTime) ScanRow(row database.RowInterface) (interface{}, error) {
	var value0 sql.NullString
	var value1 types.StringArray
	var value2 types.DateTime
	err := row.Scan(&value0, &value1, &value2)
	if err != nil {
		return nil, err
	}
	scanner.First = types.NilString()
	if value0.Valid {
		scanner.First = types.NewString(value0.String)
	}
	scanner.Second = value1
	scanner.Third = value2
	return *scanner, nil
}
