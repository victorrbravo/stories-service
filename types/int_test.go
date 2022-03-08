package types_test

import (
	"testing"

	"stories/types"
)

func TestScanInt(t *testing.T) {
	stringTests := []struct {
		name         string
		value        interface{}
		spectedValue interface{}
		spectedError bool
	}{
		// Valid Scan
		{name: "[1]: Simple int", value: int(10), spectedValue: int64(10)},
		{name: "[2]: Simple int", value: int32(20), spectedValue: int64(20)},
		{name: "[3]: Simple int", value: int64(40), spectedValue: int64(40)},
		// Invalid Scan
		{name: "[4]: Simple int", value: "1", spectedError: true},
		{name: "[5]: Simple int", value: "2", spectedError: true},
	}
	// using tt.name from the case to use it as the `t.Run` test name
	for _, tt := range stringTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			var str = new(types.Int)
			var err = str.Scan(tt.value)
			// Error quest
			if err != nil {
				if !tt.spectedError {
					t.Errorf("Error Scan | Name: %s, error: %v \n", tt.name, err)
				}
				return // On error next element
			}
			if tt.spectedError {
				t.Errorf("Error Not Found | Name: %s, value: %v \n", tt.name, tt.value)
				return
			}
			// Get Value
			// Verifica Valor
			value, err := str.Value()
			if err != nil {
				t.Errorf("Error get value | Name: %s, err %v, value: %v \n", tt.name, err, tt.value)
				return
			}

			var vl = value
			if vl != tt.spectedValue {
				t.Errorf("Error value not valid | Name: %s, spected %v, value: %v \n", tt.name, tt.spectedValue, vl)
				return
			}

		})

	}

}

func TestMarshallInt(t *testing.T) {
	stringTests := []struct {
		name    string
		value   types.Int
		spected string
	}{
		// Valid Marshal
		{name: "[1]: Simple int", value: types.NewInt(20), spected: "20"},
		{name: "[2]: Simple int", value: types.NewInt(60), spected: "60"},
		{name: "[3]: Simple int", value: types.NewInt(1200), spected: "1200"},
	}
	// using tt.name from the case to use it as the `t.Run` test name
	for _, tt := range stringTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.value.MarshalJSON()
			if err != nil {
				t.Errorf("Error MarshalJSON | Name: %s, error: %v \n", tt.name, err)
				return
			}
			if string(value) != string(tt.spected) {
				t.Errorf("Error value not valid | Name: %s, spected %v, value: %v \n", tt.name, tt.spected, value)
				return
			}
		})

	}

}

func TestUnmarshalJSONInt(t *testing.T) {
	stringTests := []struct {
		name         string
		value        string
		spected      interface{}
		spectedError bool
	}{
		// Valid Bytes Content
		{name: "[1]: Simple int", value: "20", spected: int64(20)},
		{name: "[2]: Simple int", value: "60", spected: int64(60)},
		{name: "[3]: Simple int", value: "1200", spected: int64(1200)},
		{name: "[3]: Simple int", value: "null", spected: nil},
		// Invalid Value
		{name: "[4]: Simple int", value: "a", spectedError: true},
		{name: "[5]: Simple int", value: `b`, spectedError: true},
		{name: "[6]: Simple int", value: `c`, spectedError: true},
	}

	// using tt.name from the case to use it as the `t.Run` test name

	for _, tt := range stringTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			var parseElement = new(types.Int)
			var err = parseElement.UnmarshalJSON([]byte(tt.value))
			// Error quest
			if err != nil {
				if !tt.spectedError {
					t.Errorf("Error UnmarshalJSON | Name: %s, error: %v \n", tt.name, err)
				}
				return // On error next element
			}
			if tt.spectedError {
				t.Errorf("Error Not Found | Name: %s, value: %v \n", tt.name, tt.value)
				return
			}
			// Get Value
			value, err := parseElement.Value()
			if err != nil {
				t.Errorf("Error get Value | Name: %s, error: %v \n", tt.name, err)
				return
			}
			if value != tt.spected {
				t.Errorf("Error value not valid | Name: %s, spected %v, value: %v \n", tt.name, tt.spected, value)
			}
		})

	}

}
