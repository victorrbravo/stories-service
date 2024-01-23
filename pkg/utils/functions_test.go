package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatYaml(t *testing.T) {

	input0 := `
action: update_user pk_user: 42 firstname: Hecmair
`
	value := FormatYaml(input0)
	fmt.Println("value", value)

	input1 := `
action:update_user 
pk_user: 42 firstname: Hecmair
`
	value1 := FormatYaml(input1)
	fmt.Println("value1", value1)
	input2 := `
action: update_user 
pk_user: 42 
firstname: Hecmair
`
	value2 := FormatYaml(input2)
	fmt.Println("value2", value2)

	assert.Equal(t, value, value1, "value and value1 are different")
	assert.Equal(t, value1, value2, "value1 and value2 are different")
	assert.Equal(t, value, value2, "value and value2 are different")

}
