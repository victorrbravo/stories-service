package types_test

import (
	"testing"

	"stories/types"
)

func TestHelpersData(t *testing.T) {
	stringTests := []struct {
		name         string
		value        interface{}
		spectedError bool
		spectedValue interface{}
	}{
		// Valid Scan
		{name: "[1]: Simple Helpers, bool - ", value: types.AsBool(types.NewBool(true)), spectedValue: true},
		{name: "[2]: Simple Helpers, int - ", value: types.AsInt(types.NewInt(22)), spectedValue: int(22)},
		{name: "[3]: Simple Helpers, string - ", value: types.AsString(types.NewString("hello")), spectedValue: "hello"},
		{name: "[4]: Simple Helpers, IsNill - ", value: types.IsNil(types.NewBool(true)), spectedValue: false},
	}
	// using tt.name from the case to use it as the `t.Run` test name
	for _, tt := range stringTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.spectedValue {
				t.Errorf("Error value not valid | Name: %s, spected %v, value: %v \n", tt.name, tt.spectedValue, tt.value)
				return
			}
		})
	}
}
