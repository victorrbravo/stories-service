package input

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"stories/actions/payment"
	"stories/actions/workflow"
	"stories/database"
	pg "stories/generated/inputs"
)

type ProcessorFunc func(string, string, *database.IQuery) error

type GetAction struct {
	Action string `yaml:"action"`
}

func ProcessInput(action []byte, query *database.IQuery) error {
	var getAction GetAction
	var value any
	err := yaml.Unmarshal(action, &getAction)
	if err != nil {
		return err
	}
	order := string(action)
	value, err = Processor(getAction.Action, order, query, pg.MainProcessor, payment.Processor,
		workflow.Processor)
	if err != nil {
		return err
	}
	fmt.Println("value", value)
	return nil
}

func Processor(action string, order string, query *database.IQuery, processorFuncs ...ProcessorFunc) (any, error) {
	var err error
	var value any
	for _, p := range processorFuncs {
		err = p(action, order, query)
		if err != nil {
			return value, err
		}
	}
	return value, nil
}
