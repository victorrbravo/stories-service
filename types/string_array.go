package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"stories/constants"
)

type state int

type StringArray []string

const (
	initial  state = 1 << 0
	reading        = 1 << 1
	verbose        = 1 << 2
	escaping       = 1 << 3
)

func NewStringArray(array []string) StringArray {
	return array
}

func (stringArray *StringArray) fromString(value string) error {
	var buffer strings.Builder
	var currentState state = initial
	var chr rune
	// This is valid for SQL and JSON
	if strings.ToLower(value) == constants.JSONNil {
		*stringArray = []string{}
		return nil
	}
	*stringArray = make([]string, 0)
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
				break
			case '}':
				if len(buffer.String()) > 0 {
					*stringArray = append(*stringArray, buffer.String())
				}
				return nil
			case ',':
				if currentState == reading {
					*stringArray = append(*stringArray, buffer.String())
					buffer.Reset()
				} else {
					return fmt.Errorf("unexpected `%v' found while scanning", chr)
				}
				break
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
		*stringArray = array
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (stringArray StringArray) MarshalJSON() ([]byte, error) {
	var err error
	var value []byte
	var elements []string = make([]string, 0)
	if stringArray == nil {
		return []byte(constants.JSONArrayEmpty), nil
	}
	for _, element := range stringArray {
		elements = append(elements, strings.ReplaceAll(element, "'", ""))
	}
	value, err = json.Marshal(elements)
	return value, err
}

func (stringArray *StringArray) Scan(src interface{}) error {
	switch value := src.(type) {
	case nil:
		*stringArray = nil
		break
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
	var elements []string = make([]string, len(stringArray))
	var replacer *strings.Replacer
	replacer = strings.NewReplacer("\\", "\\\\", "\"", "\\\"")
	for index, element := range stringArray {
		elements[index] = fmt.Sprintf("\"%s\"", replacer.Replace(element))
	}
	return fmt.Sprintf("{%s}", strings.Join(elements, ",")), nil
}

func (stringArray StringArray) String() string {
	if stringArray.isNil() {
		return constants.DatabaseArrayEmpty
	}
	var elements []string = make([]string, len(stringArray))
	var replacer *strings.Replacer
	replacer = strings.NewReplacer("\\", "\\\\", "\"", "\\\"")
	for index, element := range stringArray {
		elements[index] = fmt.Sprintf("\"%s\"", replacer.Replace(element))
	}
	return fmt.Sprintf("%s", strings.Join(elements, ","))
}

func (stringArray StringArray) isNil() bool {
	return stringArray == nil
}

func (stringArray StringArray) ToInterfaceArray() []interface{} {
	var interfaces []interface{}
	interfaces = make([]interface{}, len(stringArray))
	for i, element := range stringArray {
		interfaces[i] = element
	}
	return interfaces
}

func (stringArray *StringArray) FromInterfaceArray(array []interface{}) {
	*stringArray = make(StringArray, len(array))
	for i, item := range array {
		(*stringArray)[i] = item.(string)
	}
}

func (stringArray StringArray) Includes(value string) bool {
	if stringArray.isNil() {
		return false
	}

	for _, entry := range stringArray {
		if entry == value {
			return true
		}
	}

	return false
}
