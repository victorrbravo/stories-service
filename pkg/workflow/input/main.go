package input

import (
	"gopkg.in/yaml.v3"
	"log"
	"stories/actions/stories"
	"stories/constants"
	"stories/database"
	pg "stories/generated/inputs"
	"stories/pkg/utils"
	"strings"
)

type ProcessorFunc func(string, string, *database.IQuery) (string, error)

type GetAction struct {
	Action string `yaml:"action"`
}

func ProcessInput(action []byte, query *database.IQuery) (string, error) {
	var getAction GetAction
	var jsonResponse = constants.EmptyJson
	var value string
	yamlAction := utils.FormatYaml(string(action))
	log.Println("yamlAction  action:", string(action))
	err := yaml.Unmarshal([]byte(yamlAction), &getAction)
	log.Println("getAction err:", err)
	if err != nil {
		return jsonResponse, err
	}
	order := string(action)
	log.Println("pass action order:", order)
	value, err = Processor(getAction.Action, order, query, pg.MainProcessor, stories.FlowProcessor,
		stories.WorkflowProcessor)
	if err != nil {
		return jsonResponse, err
	}
	jsonResponse = value
	return value, nil
}

// FIXME: validate uniques responses
func Processor(action string, order string, query *database.IQuery, processorFuncs ...ProcessorFunc) (string, error) {
	var err error
	var value string
	for _, p := range processorFuncs {
		value, err = p(action, order, query)
		if err != nil {
			return value, err
		}
		if len(value) != 0 && strings.Compare(value, constants.EmptyJson) != 0 {
			break
		}
	}
	return value, nil
}
