package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/victorrbravo/stories-service/constants"
)

type state int

type StringArray struct {
	notNil bool
	Array  []string
}

var NilStringArray = StringArray{notNil: false, Array: nil}

const (
	initial  state = 1 << 0
	reading        = 1 << 1
	verbose        = 1 << 2
	escaping       = 1 << 3
)

func NewStringArray(array []string) StringArray {
	return StringArray{true, array}
}

func (stringArray *StringArray) fromString(value string) error {
	var buffer strings.Builder
	var currentState state = initial
	var chr rune
	// This is valid for SQL and JSON
	if strings.ToLower(value) == constants.JSONNil {
		*stringArray = NilStringArray
		return nil
	}
	stringArray.Array = []string{}
	stringArray.notNil = true
	for _, chr = range value {
		if currentState&escaping == escaping {
			currentState = currentState ^ escaping
			buffer.WriteRune(chr)
		} else if chr == '\\' {
			currentState = currentState ^ escaping
		} else if currentState&verbose == verbose {
			if chr == '"' {
				currentState = currentState ^ verbose
			} else {
				buffer.WriteRune(chr)
			}
		} else if chr == '"' {
			currentState = currentState | verbose
		} else {
			switch chr {
			case ' ', '\t', '\n':
				break
			case '{':
				if currentState != initial {
					return fmt.Errorf("expecting the `{' character but found `%c'", chr)
				}
				currentState = reading
			case '}':
				if len(buffer.String()) > 0 {
					stringArray.Array = append(stringArray.Array, buffer.String())
				}
				return nil
			case ',':
				if currentState == reading {
					stringArray.Array = append(stringArray.Array, buffer.String())
					buffer.Reset()
				} else {
					return fmt.Errorf("unexpected `%v' found while scanning", chr)
				}
			default:
				if currentState == reading {
					buffer.WriteRune(chr)
				} else {
					return fmt.Errorf("expecting a control carachter `}', `{' or `,'")
				}
			}
		}
	}
	return fmt.Errorf("end of input reached but `}' was never found")
}

func (stringArray *StringArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value []string
	value = stringArray.Array
	if err := unmarshal(&value); err != nil {
		return err
	}
	*stringArray = StringArray{true, value}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (stringArray *StringArray) UnmarshalJSON(src []byte) error {
	strValue := string(src)
	var array []string
	if strValue != constants.JSONNil {
		var err error
		err = json.Unmarshal(src, &array)
		if err != nil {
			return err
		}
		stringArray.Array = array
		stringArray.notNil = true
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (stringArray StringArray) MarshalJSON() ([]byte, error) {
	var err error
	var value []byte
	var elements []string = make([]string, 0)
	if stringArray.isNil() {
		return []byte(constants.JSONArrayEmpty), nil
	}
	for _, element := range stringArray.Array {
		elements = append(elements, strings.ReplaceAll(element, "'", ""))
	}
	value, err = json.Marshal(elements)
	return value, err
}

func (stringArray *StringArray) Scan(src interface{}) error {
	switch value := src.(type) {
	case nil:
		*stringArray = NilStringArray
	case string:
		return stringArray.fromString(value)
	case []byte:
		return stringArray.fromString(string(value))
	default:
		return fmt.Errorf(constants.CanntReadAStringArrayFromNonString)
	}
	return nil
}

func (stringArray StringArray) Value() (driver.Value, error) {
	if stringArray.isNil() {
		return nil, nil
	}
	var elements []string = make([]string, len(stringArray.Array))
	var replacer *strings.Replacer
	replacer = strings.NewReplacer("\\", "\\\\", "\"", "\\\"")
	for index, element := range stringArray.Array {
		elements[index] = fmt.Sprintf("\"%s\"", replacer.Replace(element))
	}
	return fmt.Sprintf("{%s}", strings.Join(elements, ",")), nil
}

func (stringArray StringArray) String() string {
	if stringArray.isNil() {
		return constants.DatabaseArrayEmpty
	}
	var elements []string = make([]string, len(stringArray.Array))
	var replacer *strings.Replacer
	replacer = strings.NewReplacer("\\", "\\\\", "\"", "\\\"")
	for index, element := range stringArray.Array {
		elements[index] = fmt.Sprintf("\"%s\"", replacer.Replace(element))
	}
	return fmt.Sprintf("%s", strings.Join(elements, ","))
}

func (stringArray StringArray) isNil() bool {
	return !stringArray.notNil
}

func (stringArray StringArray) ToInterfaceArray() []interface{} {
	var interfaces []interface{}
	interfaces = make([]interface{}, len(stringArray.Array))
	for i, element := range stringArray.Array {
		interfaces[i] = element
	}
	return interfaces
}

func (stringArray *StringArray) FromInterfaceArray(array []interface{}) {
	stringArray.Array = make([]string, len(array))
	stringArray.notNil = true
	for i, item := range array {
		stringArray.Array[i] = item.(string)
	}
}

func (stringArray StringArray) Includes(value string) bool {
	if stringArray.isNil() {
		return false
	}

	for _, entry := range stringArray.Array {
		if entry == value {
			return true
		}
	}

	return false
}
