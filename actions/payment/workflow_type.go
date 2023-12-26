package payment

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"stories/constants"
	"stories/database"
	"stories/generated/models"
	"stories/pkg/workflow/common"
)

// String a boolean type that understands jsons `null'
type WorkflowType struct {
	notNil bool
	value  string
	wf     *common.Workflow[models.TaTracking]
	Query  *database.IQuery
}

func (wft WorkflowType) Workflow() *common.Workflow[models.TaTracking] {
	if wft.wf == nil {
		log.Print("worflow is nil for WorkflowType")
	}
	return wft.wf
}

// Value implements the sql.Valuer interface
func (str WorkflowType) Value() (driver.Value, error) {
	if !str.notNil {
		return nil, nil
	}
	return str.value, nil
}

// Scan implements the sql.Scanner interface
func (str *WorkflowType) Scan(src interface{}) error {
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
func (str *WorkflowType) UnmarshalJSON(src []byte) error {
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

func (wft *WorkflowType) LoadWorkflow() (*common.Workflow[models.TaTracking], error) {
	ctx := context.Background()
	dataFile, err := GetJsonDatafile(wft.Query)
	if err != nil {
		return nil, err
	}
	wft.wf = new(common.Workflow[models.TaTracking])
	err = wft.wf.Load(ctx, dataFile)
	if err != nil {
		return wft.wf, err
	}
	return wft.wf, nil
}

func (wft *WorkflowType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	var err error
	ctx := context.Background()
	value = wft.value
	if err := unmarshal(&value); err != nil {
		return err
	}
	if wft.wf == nil {
		wft.wf, err = wft.LoadWorkflow()
		if err != nil {
			return err
		}
	}
	statuses, err := wft.wf.Statuses(ctx, false)
	if err != nil {
		return err
	}
	index := slices.Index(statuses, value)
	if index == -1 {
		return fmt.Errorf("error: value '%s' is not found in list '%v'", value, statuses)
	}
	*wft = NewWorkflowType(value)
	return nil
}

func (str WorkflowType) MarshalYAML() (interface{}, error) {
	return str.value, nil
}

// MarshalJSON implements the json.Marshaler interface
func (str WorkflowType) MarshalJSON() ([]byte, error) {
	if !str.notNil {
		return json.Marshal(nil)
	}
	return json.Marshal(str.value)
}

func (str WorkflowType) isNil() bool {
	return !str.notNil
}

func (str WorkflowType) get() string {
	return str.value
}

// NewString builds a nullable string from a native string
func NewWorkflowType(str string) WorkflowType {
	return WorkflowType{value: str, notNil: true}
}

func NilWorkflowType() WorkflowType {
	return WorkflowType{value: "", notNil: false}
}

func (str WorkflowType) String() string {
	if str.notNil {
		return str.value
	} else {
		return "nil"
	}
}
