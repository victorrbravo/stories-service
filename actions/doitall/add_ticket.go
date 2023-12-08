package doitall

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"stories/actions/constants"
	"stories/database"
	"stories/generated/models"
	"stories/types"
)

type Ticket struct {
	Name    types.String `yaml:"name"`
	Summary types.String `yaml:"summary"`
}

type ChangeStatus struct {
	Ticket types.Int    `yaml:"ticket"`
	Status TicketStatus `yaml:"status"`
}

type DeleteTicket struct {
	Ticket types.Int `yaml:"ticket"`
}

type TicketStatus string

const (
	ToDo     TicketStatus = "ToDo"
	Progress TicketStatus = "Progress"
	Done     TicketStatus = "Done"
)

func (cs TicketStatus) IsValid() bool {
	switch cs {
	case ToDo, Progress, Done:
		return true
	}
	return false
}

func (cs *TicketStatus) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	*cs = TicketStatus(value)
	if !cs.IsValid() {
		err := errors.New(fmt.Sprintf("'%s' is not valid. Valid Values ToDo, Progress, Done", value))
		return err
	}
	return nil
}

func Processor(action string, order string, query *database.IQuery) error {
	var err error
	switch action {
	case "add_ticket":
		err = AddTicket(order, query)
	case "change_status":
		err = ChangeTicketStatus(order, query)
	case "delete_ticket":
		err = DoDeleteTicket(order, query)
	}
	if err != nil {
		return err
	}
	return nil
}

func UnmarshallTicket(order []byte) (Ticket, error) {
	var ticket Ticket
	var err error
	err = yaml.Unmarshal(order, &ticket)
	if err != nil {
		return ticket, err
	}

	return ticket, nil
}

func UnmarshallChangeStatus(order []byte) (ChangeStatus, error) {
	var change ChangeStatus
	var err error
	err = yaml.Unmarshal(order, &change)
	if err != nil {
		return change, err
	}
	return change, nil
}

func AddTicket(order string, query *database.IQuery) error {
	handler, err := UnmarshallTicket([]byte(order))
	if err != nil {
		return err
	}

	toDo := types.NewString(constants.ToDo)
	images := types.NewStringArray([]string{})
	taStory, err := models.CreateTaStory(*query,
		handler.Name,
		toDo,
		handler.Summary,
	)
	if err != nil {
		return err
	}
	models.CreateTaTracking(*query,
		taStory.PkStory,
		toDo,
		images,
	)
	return nil
}

func ChangeTicketStatus(order string, query *database.IQuery) error {
	handler, err := UnmarshallChangeStatus([]byte(order))
	if err != nil {
		return err
	}

	taStory, err := models.GetTaStoryByPkStory(*query,
		handler.Ticket)
	if err != nil {
		return err
	}

	images := types.NewStringArray([]string{})
	status := types.NewString(string(handler.Status))
	taStory.Update(*query,
		types.NilString(),
		status,
		types.NilString())

	models.CreateTaTracking(*query,
		taStory.PkStory,
		status,
		images,
	)

	return nil
}

func DoDeleteTicket(order string, query *database.IQuery) error {
	handler, err := UnmarshallChangeStatus([]byte(order))
	if err != nil {
		return err
	}
	err = models.DeleteTaTrackingByFkStory(*query, handler.Ticket)
	if err != nil {
		return err
	}
	taStory, err := models.GetTaStoryByPkStory(*query, handler.Ticket)
	if err != nil {
		return err
	}
	err = taStory.Delete(*query)
	if err != nil {
		return err
	}
	return nil
}
