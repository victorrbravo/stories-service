package types

import "fmt"

// Float32OneDecimal a special float value that marshals json to a single decimal place
type Float32OneDecimal float32

// MarshalJSON implements the json.Marshaler interface
func (value Float32OneDecimal) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.1f", value)), nil
}
