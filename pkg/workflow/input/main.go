package input

import (
	"gopkg.in/yaml.v3"
	"stories/database"
	pg "stories/generated/inputs"
)

type GetAction struct {
	Action string `yaml:"action"`
}

func ProcessInput(action []byte, query *database.IQuery) error {
	var getAction GetAction
	err := yaml.Unmarshal(action, &getAction)
	if err != nil {
		return err
	}
	order := string(action)
	err = pg.ChooseProcessor(getAction.Action, order, query)
	if err != nil {
		return err
	}
	return nil
}
