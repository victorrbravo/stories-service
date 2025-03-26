package stories

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"stories/constants"
	"stories/database"
	"stories/pkg/workflow/common"
)

// String a boolean type that understands jsons `null'
type WorkflowType[T database.Scanner] struct {
	notNil bool
	value  string
	wf     *common.Workflow[T]
	Query  *database.IQuery
	PkFlow int
}

func (wft WorkflowType[T]) Workflow() *common.Workflow[T] {
	if wft.wf == nil {
		log.Print("worflow is nil for WorkflowType")
	}
	return wft.wf
}

// Value implements the sql.Valuer interface
func (str WorkflowType[T]) Value() (driver.Value, error) {
	if !str.notNil {
		return nil, nil
	}
	return str.value, nil
}

// Scan implements the sql.Scanner interface
func (str *WorkflowType[T]) Scan(src interface{}) error {
	switch value := src.(type) {
	case nil:
		str.notNil = false
		str.value = ""
		break
	case string:
		str.notNil = true
		str.value = value
		break
	case []byte:
		str.notNil = true
		str.value = string(value)
	default:
		return fmt.Errorf(constants.InvalidValueForString, value)
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (str *WorkflowType[T]) UnmarshalJSON(src []byte) error {
	strValue := string(src)
	if strValue != constants.JSONNil {
		var err error
		err = json.Unmarshal(src, &str.value)
		if err != nil {
			return err
		}
		str.notNil = true
	}
	return nil
}

func (wft *WorkflowType[T]) LoadWorkflow(pkFlow int) (*common.Workflow[T], error) {
	ctx := context.Background()
	dataFile, err := GetJsonDatafile(pkFlow, wft.Query)
	if err != nil {
		return nil, err
	}
	wft.wf = new(common.Workflow[T])
	err = wft.wf.Load(ctx, dataFile)
	if err != nil {
		return wft.wf, err
	}
	return wft.wf, nil
}

func (wft *WorkflowType[T]) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	var err error
	wft.PkFlow = 1
	ctx := context.Background()
	value = wft.value
	if err := unmarshal(&value); err != nil {
		return err
	}
	if wft.wf == nil {
		wft.wf, err = wft.LoadWorkflow(wft.PkFlow)
		if err != nil {
			return err
		}
	}
	statuses, err := wft.wf.Statuses(ctx, false, true)
	if err != nil {
		return err
	}
	index := slices.Index(statuses, value)
	if index == -1 {
		return fmt.Errorf("error: value '%s' is not found in list '%v'", value, statuses)
	}
	*wft = NewWorkflowType[T](value)
	return nil
}

func (str WorkflowType[T]) MarshalYAML() (interface{}, error) {
	return str.value, nil
}

// MarshalJSON implements the json.Marshaler interface
func (str WorkflowType[T]) MarshalJSON() ([]byte, error) {
	if !str.notNil {
		return json.Marshal(nil)
	}
	return json.Marshal(str.value)
}

func (str WorkflowType[T]) isNil() bool {
	return !str.notNil
}

func (str WorkflowType[T]) get() string {
	return str.value
}

// NewString builds a nullable string from a native string
func NewWorkflowType[T database.Scanner](str string) WorkflowType[T] {
	return WorkflowType[T]{value: str, notNil: true}
}

func NilWorkflowType[T database.Scanner]() WorkflowType[T] {
	return WorkflowType[T]{value: "", notNil: false}
}

func (str WorkflowType[T]) String() string {
	if str.notNil {
		return str.value
	} else {
		return "nil"
	}
}
