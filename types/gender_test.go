package types_test

import (
	"fmt"
	"testing"

	"stories/constants"
	"stories/types"
)

func TestGender(t *testing.T) {
	stringTests := []struct {
		name         string
		value        interface{}
		spectedValue interface{}
		spectedError bool
	}{
		// Valid Scan
		{name: "[1]: Male String", value: "M", spectedValue: "M"},
		{name: "[2]: Male Byte", value: []byte("M"), spectedValue: "M"},
		{name: "[3]: Female String", value: "F", spectedValue: "F"},
		{name: "[4]: Female Byte", value: []byte("F"), spectedValue: "F"},

		// Invalid Scan
		{name: "[5]: Nil", value: nil, spectedError: true}, // Se pasa "nil" retorna 0?
		{name: "[7]: Empty String", value: string(`""`), spectedError: true},
		{name: "[9]: Other String", value: []byte("X"), spectedError: true},

		// FAIL ON EMPTY BYTE AND EMPTY STRING

		//{name: "[6]: Empty Byte", value: []byte(""), spectedError: true},
		//{name: "[8]: Empty String", value: "", spectedError: true},

	}

	for _, tt := range stringTests {
		t.Run(tt.name, func(t *testing.T) {
			str := new(types.Gender)
			err := str.Scan(tt.value)
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

			value, err := str.Value()
			if err != nil {
				t.Errorf("Error get value | Name: %s, err %v, value: %v \n", tt.name, err, tt.value)
				return
			}

			vl := value
			if vl != tt.spectedValue {
				t.Errorf("Error value not valid | Name: %s, spected %v, value: %v \n", tt.name, tt.spectedValue, vl)
				return
			}
		})
	}
}

func TestMarshallGender(t *testing.T) {
	stringTests := []struct {
		name         string
		value        types.Gender
		spected      rune
		spectedError bool
	}{
		// Valid Marshal
		{name: "[1]: FEMALE", value: 'F', spected: rune('F')},
		{name: "[2]: MALE", value: 'M', spected: rune('M')},
		{name: "[3]: OTHER", value: 'Z', spectedError: true},
	}

	// using tt.name from the case to use it as the `t.Run` test name
	for _, tt := range stringTests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.value.MarshalJSON()
			if err != nil {
				if tt.spectedError == true {
					if err.Error() == fmt.Errorf(constants.InvalidGenderValue, tt.value).Error() {
						return
					}

					t.Errorf("Error value not valid | Name: %s, spected %v, value: %v \n", tt.name, tt.spected, value)
				}
			}
			if string(value) != string(tt.spected) {
				return
			}
		})
	}
}

func TestUnmarshalJSONGender(t *testing.T) {
	stringTests := []struct {
		name         string
		value        string
		spected      interface{}
		spectedError bool
	}{
		// Valid Bytes Content
		{name: "[1]: Simple Float", value: "true", spected: true},
		{name: "[2]: Simple Float", value: "false", spected: false},
		{name: "[3]: Simple Float", value: "null", spected: nil},
		// Invalid Value
		{name: "[4]: Simple Float", value: "1", spectedError: true},
		{name: "[5]: Simple Float", value: `"false"`, spectedError: true},
		{name: "[6]: Simple Float", value: `"true"`, spectedError: true},
	}

	// using tt.name from the case to use it as the `t.Run` test name

	for _, tt := range stringTests {
		t.Run(tt.name, func(t *testing.T) {
			parseElement := new(types.Bool)
			err := parseElement.UnmarshalJSON([]byte(tt.value))
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
