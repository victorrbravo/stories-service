package scanners

import (
	"database/sql"
	"stories/database"
	"stories/types"
)

type TupleTwo struct {
	First  types.String
	Second types.String
}

func (scanner *TupleTwo) ScanRow(row database.RowInterface) (interface{}, error) {
	var value0, value1 sql.NullString
	var err error
	err = row.Scan(&value0, &value1)
	if err != nil {
		return nil, err
	}
	scanner.First = types.NilString()
	scanner.Second = types.NilString()

	if value0.Valid {
		scanner.First = types.NewString(value0.String)
	}
	if value1.Valid {
		scanner.Second = types.NewString(value1.String)
	}
	return scanner, nil
}
