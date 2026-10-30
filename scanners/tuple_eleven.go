package scanners

import (
	"database/sql"
	"github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/types"
)

type TupleEleven struct {
	First    types.String
	Second   types.String
	Third    types.String
	Four     types.String
	ID       types.Int
	Fifth    types.StringArray
	Sixth    types.StringArray
	Seventh  types.Int
	Eight    types.Int
	Ninth    types.Coordinates
	Tenth    types.String
	Eleventh types.StringArray
}

func (scanner *TupleEleven) ScanRow(row database.RowInterface) (interface{}, error) {
	var value0, value1, value2, value3 sql.NullString
	var valueID types.Int
	var value5, value6 types.StringArray
	var value7 types.Int
	var value8 types.Int
	var value9 types.Coordinates
	var value10 types.String
	var value11 types.StringArray

	var err error
	err = row.Scan(&value0, &value1, &value2, &value3, &valueID, &value5, &value6, &value7,
		&value8,
		&value9,
		&value10,
		&value11)

	if err != nil {
		return nil, err
	}
	scanner.First = types.NilString()
	scanner.Second = types.NilString()
	scanner.Third = types.NilString()
	scanner.Four = types.NilString()
	scanner.ID = types.NilInt
	scanner.Seventh = types.NilInt
	scanner.Eight = types.NilInt
	scanner.Ninth = types.NilCoordinates
	scanner.Tenth = types.NilString()
	scanner.Eleventh = types.NilStringArray

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
	scanner.Seventh = value7
	scanner.Eight = value8
	scanner.Ninth = value9
	scanner.Tenth = value10
	scanner.Eleventh = value11

	return scanner, nil
}
