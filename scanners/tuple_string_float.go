package scanners

import (
	"database/sql"

	"github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/types"
)

// TupleStringFloat scans a (text, numeric) row, i.e. a grouping key plus its
// aggregated value.
type TupleStringFloat struct {
	First  types.String
	Second types.Float
}

// ScanRow implements the database.Scanner interface. It returns a copy of the
// scanner, not the scanner itself, so the rows collected by QueryMany do not
// all alias the same struct.
func (scanner *TupleStringFloat) ScanRow(row database.RowInterface) (interface{}, error) {
	var value0 sql.NullString
	var value1 sql.NullFloat64
	var err error
	err = row.Scan(&value0, &value1)
	if err != nil {
		return nil, err
	}
	scanner.First = types.NilString()
	scanner.Second = types.NilFloat()

	if value0.Valid {
		scanner.First = types.NewString(value0.String)
	}
	if value1.Valid {
		scanner.Second = types.NewFloat(value1.Float64)
	}
	return *scanner, nil
}
