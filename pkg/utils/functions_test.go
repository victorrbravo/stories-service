package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestFormatYamlFlowStep(t *testing.T) {

	input1 := `
action: update_user step_id: 123 task:	name: victor user_id: 12
`
	value1 := FormatYaml(input1)
	fmt.Println("value1\n", value1)

	input := `
 action: add_flow_step ticket: 44 status: Shopping mall images: ABC897B7-71DA-4C08-B7E8-B7F2F071D445, amount: 3500.00 summary: In CC Cacique media_date: 2024-03-24T18:05:26Z type: step
`
	value := FormatYaml(input)
	fmt.Println("value", value)

}

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

	input3 := `
action: update_user images: 0123-1567-89ab-cdef, 0123-2562-89ab-cdef, 3024-3567-89ab-cdef, 5023-3567-89ab-cdef  name: Hecmair 
`
	expected3 := `
action: update_user 
images: 
 - 0123-1567-89ab-cdef 
 - 0123-2562-89ab-cdef 
 - 3024-3567-89ab-cdef 
 - 5023-3567-89ab-cdef  
name: Hecmair 
`
	value3 := FormatYaml(input3)
	assert.Equal(t, value3, expected3, "value3 and expected3 are different")

}

type InputTest1 struct {
	Name   string   `yaml:"name"`
	Action string   `yaml:"action"`
	Images []string `yaml:"images"`
	Roles  []string `yaml:"roles"`
}

func TestFormatYamlArray(t *testing.T) {

	input3 := `
action: update_user images: 0123-1567-89ab-cdef, 0123-2562-89ab-cdef, 3024-3567-89ab-cdef, 5023-3567-89ab-cdef  roles: state 1,state 2, state 3 name: Hecmair 
`
	expected3 := `action: update_user
images: 
 - 0123-1567-89ab-cdef
 - 0123-2562-89ab-cdef
 - 3024-3567-89ab-cdef
 - 5023-3567-89ab-cdef
roles: 
 - state 1
 - state 2
 - state 3
name: Hecmair
`
	value3 := FormatYaml(input3)
	order := []byte(value3)
	var inputTest1 InputTest1
	err := yaml.Unmarshal(order, &inputTest1)
	if err != nil {
		t.Errorf("error %v", err)
	}
	fmt.Printf("inputTest1: %+v\n", inputTest1)
	assert.Equal(t, value3, expected3, "value3 and expected3 are different")

}

func TestFormatYamlArray2(t *testing.T) {
	input3 := `
action: update_user images: 0123-1567-89ab-cdef, roles: state 1,state 2, state 3 name: Hecmair 
`
	expected3 := `action: update_user
images: 
 - 0123-1567-89ab-cdef
roles: 
 - state 1
 - state 2
 - state 3
name: Hecmair
`
	value3 := FormatYaml(input3)
	order := []byte(value3)
	var inputTest1 InputTest1
	err := yaml.Unmarshal(order, &inputTest1)
	if err != nil {
		t.Errorf("error %v", err)
	}
	fmt.Printf("inputTest1: %+v\n", inputTest1)
	assert.Equal(t, value3, expected3, "value3 and expected3 are different")

}
