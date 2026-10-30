package scanners

import (
	"database/sql"
	"github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/types"
)

type TupleThree struct {
	First  types.String
	Second types.String
	Third  types.String
}

func (scanner *TupleThree) ScanRow(row database.RowInterface) (interface{}, error) {
	var value0, value1, value3 sql.NullString
	var err error
	err = row.Scan(&value0, &value1, &value3)
	if err != nil {
		return nil, err
	}
	scanner.First = types.NilString()
	scanner.Second = types.NilString()
	scanner.Third = types.NilString()

	if value0.Valid {
		scanner.First = types.NewString(value0.String)
	}
	if value1.Valid {
		scanner.Second = types.NewString(value1.String)
	}
	if value3.Valid {
		scanner.Third = types.NewString(value3.String)
	}
	return scanner, nil
}
