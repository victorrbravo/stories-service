package scanners

import (
	"stories/database"
	"stories/types"
)

type UUID struct {
	Value types.UUID
}

func (scanner *UUID) ScanRow(row database.RowInterface) (interface{}, error) {
	var err error
	var uuidValue types.UUID
	err = row.Scan(&uuidValue)
	if err != nil {
		return nil, err
	}
	scanner.Value = uuidValue
	return *scanner, nil
}
