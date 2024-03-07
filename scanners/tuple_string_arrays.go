package scanners

import (
	"stories/database"
	"stories/types"
)

type TupleArrays struct {
	First  types.StringArray
	Second types.StringArray
	Third  types.StringArray
}

func (scanner *TupleArrays) ScanRow(row database.RowInterface) (interface{}, error) {
	var value0, value1, value2 types.StringArray
	var err error
	err = row.Scan(&value0, &value1, &value2)
	if err != nil {
		return nil, err
	}
	scanner.First = value0
	scanner.Second = value1
	scanner.Third = value2
	return scanner, nil
}
