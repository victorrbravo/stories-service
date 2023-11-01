package input

import "testing"

const testYaml1 = `
action: add
summary: replace function using sync function algorithm
due_date: 07/10/11
project: ct-101
`

const testYaml2 = `
action: add
fields:
  summary: replace function using sync function algorithm
  due_date: 07/10/11
  project: ct-101

`

func TestReadAction(t *testing.T) {
	err := ReadAction(testYaml1)
	if err == nil {
		t.Errorf("readAction err must be distinct of nil. err %v\n", err)
	}
	err = ReadAction(testYaml2)
	if err != nil {
		t.Errorf("ReadAction fails err %v\n", err)
	}

}
