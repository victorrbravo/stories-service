package scanners

import (
	"database/sql"
	"fmt"
	"reflect"

	"stories/database"
)

type Reflect struct {
	interfaceValue interface{}
}

func NewReflectScanner(value interface{}) Reflect {
	return Reflect{interfaceValue: value}
}

func (instance Reflect) ScanRow(row database.RowInterface) (interface{}, error) {
	var value interface{}
	var err error
	if sqlRows, ok := row.(*sql.Rows); ok {
		value, err = scanAndDiscover(instance.interfaceValue, sqlRows)
		if err != nil {
			return nil, err
		}
		return value, nil
	} else {
		return nil, fmt.Errorf("cannot scan a single row with reflection")
	}
}

func scanAndDiscover(scannerValue interface{}, rows *sql.Rows) (interface{}, error) {
	var err error
	var columnNames []string
	columnNames, err = rows.Columns()
	if err != nil {
		return nil, err
	}
	var pointers []interface{}
	var value reflect.Value
	// Create some pointers to hold the scanned values
	pointers = make([]interface{}, len(columnNames))
	// Now apply reflection to the interesting interface
	value = reflect.ValueOf(scannerValue)
	if value.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("when using the `Reflect' scanner you must initialize it with a pointer")
	}
	for index, columnName := range columnNames {
		var fieldValue reflect.Value
		fieldValue = value.Elem().FieldByName(columnName)
		if !fieldValue.IsValid() {
			return nil, fmt.Errorf("column `%s' cannot be mapped to the interface", columnName)
		}
		pointers[index] = fieldValue.Addr().Interface()
	}
	err = rows.Scan(pointers...)
	if err != nil {
		return nil, err
	}
	return scannerValue, err
}
