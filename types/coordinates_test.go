package types_test

import (
	"testing"

	"stories/types"
)

// func TestScanCoordinates(t *testing.T) {

// 	stringTests := []struct {
// 		name         string
// 		value        []byte
// 		spectedError bool
// 		spectedValue interface{}
// 	}{
// 		{value: []byte(fmt.Sprintf("%x%x", float64(1.22), float64(1.22)))},
// 	}

// 	// using tt.name from the case to use it as the `t.Run` test name
// 	for _, tt := range stringTests {
// 		// using tt.name from the case to use it as the `t.Run` test name
// 		t.Run(tt.name, func(t *testing.T) {
// 			var str = new(types.Coordinates)
// 			var err = str.Scan(tt.value)
// 			// Error quest
// 			if err != nil {
// 				if !tt.spectedError {
// 					t.Errorf("Error Scan | Name: %s, error: %v \n", tt.name, err)
// 				}
// 				return // On error next element
// 			}
// 			if tt.spectedError {
// 				t.Errorf("Error Not Found | Name: %s, value: %v \n", tt.name, tt.value)
// 				return
// 			}
// 			// Get Value
// 			value, err := str.Value()
// 			if err != nil {
// 				t.Errorf("Error get value | Name: %s, err %v, value: %v \n", tt.name, err, tt.value)
// 				return
// 			}

// 			var vl = value
// 			if vl != tt.spectedValue {
// 				t.Errorf("Error value not valid | Name: %s, spected %v, value: %v \n", tt.name, tt.spectedValue, vl)
// 				return
// 			}
// 		})
// 	}
// }

func TestMarshallCoordinates(t *testing.T) {
	stringTests := []struct {
		name    string
		value   types.Coordinates
		spected string
	}{
		// Valid Marshal
		{name: "[1]: Simple Coordinates", value: types.NewCoordinates(1.333, 1.44), spected: `{"latitude":1.333,"longitude":1.44}`},
		{name: "[2]: Simple Coordinates", value: types.NewCoordinates(4.2212, 6.2123123), spected: `{"latitude":4.2212,"longitude":6.2123123}`},
		{name: "[3]: Simple Coordinates", value: types.Coordinates{}, spected: "null"},
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
				t.Errorf("Error value not valid | Name: %s, spected %v, value: %s \n", tt.name, tt.spected, value)
				return
			}
		})
	}
}

func TestUnmarshalJSONCoordinates(t *testing.T) {
	stringTests := []struct {
		name         string
		value        string
		spected      interface{}
		spectedError bool
	}{
		// Valid Bytes Content
		{name: "[1]: Simple Coordinates", value: `{"latitude":1.333,"longitude":2.6543}`, spected: `POINT (1.333000 2.654300)`},
		{name: "[2]: Simple Coordinates", value: `{"latitude":2.6543,"longitude":1.44}`, spected: `POINT (2.654300 1.440000)`},
		// Invalid Value
		{name: "[4]: Simple Coordinates", value: `{"latitude":1.333}`, spectedError: true},
		{name: "[5]: Simple Coordinates", value: `{"longitude":1.44}`, spectedError: true},
		{name: "[6]: Simple Coordinates", value: `{}`, spectedError: true},
		// Null Elements
		{name: "[7]: Simple Coordinates", value: `{"latitude":null,"longitude":1.44}`, spected: `POINT (0.000000 1.440000)`},
		{name: "[8]: Simple Coordinates", value: `{"latitude":2.6543,"longitude":null}`, spected: `POINT (2.654300 0.000000)`},
	}

	// using tt.name from the case to use it as the `t.Run` test name

	for _, tt := range stringTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			var parseElement = new(types.Coordinates)
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
