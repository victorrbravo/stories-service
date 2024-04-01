package scanners

import (
	"database/sql"
	"stories/database"
	"stories/types"
)

type TupleSixIntArray struct {
	First  types.String
	Second types.String
	Third  types.String
	Four   types.String
	ID     types.Int
	Fifth  types.StringArray
	Sixth  types.StringArray
}

func (scanner *TupleSixIntArray) ScanRow(row database.RowInterface) (interface{}, error) {
	var value0, value1, value2, value3 sql.NullString
	var valueID types.Int
	var value5, value6 types.StringArray
	var err error
	err = row.Scan(&value0, &value1, &value2, &value3, &valueID, &value5, &value6)
	if err != nil {
		return nil, err
	}
	scanner.First = types.NilString()
	scanner.Second = types.NilString()
	scanner.Third = types.NilString()
	scanner.Four = types.NilString()
	scanner.ID = types.NilInt

	if value0.Valid {
		scanner.First = types.NewString(value0.String)
	}
	if value1.Valid {
		scanner.Second = types.NewString(value1.String)
	}
	if value2.Valid {
		scanner.Third = types.NewString(value2.String)
	}
	if value3.Valid {
		scanner.Four = types.NewString(value3.String)
	}
	scanner.ID = valueID
	scanner.Fifth = value5
	scanner.Sixth = value6
	return scanner, nil
}
