package types_test

import (
	"fmt"
	"testing"

	"github.com/victorrbravo/stories-service/constants"
	"github.com/victorrbravo/stories-service/types"
)

func TestScanString(t *testing.T) {
	stringTests := []struct {
		name       string
		value      interface{}
		spectedErr error
	}{
		{name: "", value: "text", spectedErr: nil},
		{name: "", value: 3, spectedErr: fmt.Errorf(constants.InvalidValueForString, 3)},
		{name: "", value: []byte("text"), spectedErr: nil},
		{name: "", value: []byte("null"), spectedErr: nil},
	}

	// using tt.name from the case to use it as the `t.Run` test name

	for _, tt := range stringTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			str := new(types.String)
			err := str.Scan(tt.value)
			ErrText := ""
			SpectedErrText := ""
			if err != nil {
				ErrText = err.Error()
			}
			if tt.spectedErr != nil {
				SpectedErrText = err.Error()
			}
			if ErrText != SpectedErrText {
				t.Errorf("%s got %v want %v", tt.name, err, tt.spectedErr)
				t.FailNow()
			}
		})
	}
}

func TestUnmarshalJSONString(t *testing.T) {
	stringTests := []struct {
		name         string
		value        []byte
		spected      string
		spectedError bool
		isNil        bool
	}{
		// Valid Bytes Content
		{name: "Filled Byte Array", value: []byte(`"hello text"`), spected: "hello text", spectedError: false},
		{name: "Null Text Byte Array", value: []byte(`null`), spected: "", spectedError: false, isNil: true},
		{name: "Null Text Byte Array", value: []byte(`"null"`), spected: "null", spectedError: false, isNil: false},
		// Invalid Bytes Content
		{name: "Null Byte Array", value: nil, spectedError: true},
		{name: "Null Text Byte Array", value: []byte(`nil`), spectedError: true},
		{name: "Null Text Byte Array", value: []byte(`1`), spectedError: true},
		{name: "Empty Byte Array", value: []byte{}, spectedError: true},
	}

	// using tt.name from the case to use it as the `t.Run` test name

	for _, tt := range stringTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			str := new(types.String)
			err := str.UnmarshalJSON(tt.value)
			if err != nil && !tt.spectedError {
				t.Errorf("1 %s got %v want %v", tt.name, err, tt.spected)
				return
			}
			if tt.spectedError {
				return
			}
			value, err := str.Value()
			if err != nil {
				t.Errorf("2 %s got %v want %v", tt.name, err, tt.spected)
				return
			}
			if tt.isNil && value != nil {
				t.Errorf("4 %s got %v want %v", tt.name, err, tt.spected)
			}
			if tt.isNil {
				return
			}
			if value != tt.spected {
				t.Errorf("3 %s got %v want %v", tt.name, err, tt.spected)
				return
			}
		})
	}
}

func TestMarshalJSONString(t *testing.T) {
	bytesNill := []byte{110, 117, 108, 108}
	stringTests := []struct {
		name    string
		value   types.String
		spected []byte
	}{
		// Valid Bytes Content
		{name: "Filled Text", value: types.NewString("hello text"), spected: []byte(`"hello text"`)},

		{name: "Nill Text", value: types.NilString(), spected: bytesNill},
	}

	for _, tt := range stringTests {
		// using tt.name from the case to use it as the `t.Run` test name
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.value.MarshalJSON()
			if err != nil {
				t.Errorf("%s got %v want %v", tt.name, err, tt.spected)
				return
			}

			if string(value) != string(tt.spected) {
				t.Errorf("%s got %v want %v", tt.name, string(value), string(tt.spected))
				return
			}
		})
	}
}
