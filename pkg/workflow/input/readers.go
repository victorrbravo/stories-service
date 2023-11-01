package input

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
)

func checkFields(values map[string]interface{}) error {
	_, ok := values["action"]
	if !ok {
		return errors.New("action is not present")
	}
	_, ok = values["fields"]
	if !ok {
		return errors.New("fields is not present")
	}
	return nil
}

func ReadAction(action string) error {
	// Read YAML file

	// Create a Config struct to hold the data
	data := []byte(action)
	// Unmarshal YAML data into the Config struct
	var values map[string]interface{}
	err := yaml.Unmarshal(data, &values)
	if err != nil {
		log.Fatal("Error unmarshaling YAML:", err)
	}
	err = checkFields(values)
	if err != nil {
		return err
	}
	currentAction, _ := values["action"]
	fmt.Println("currentAction", currentAction)
	fields := values["fields"].(map[string]interface{})
	for key, value := range fields {
		fmt.Println("key", key, "value", value)
	}
	return err
}
