package payment

import (
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

type AddPaymentStatus struct {
	Ticket      types.Int    `yaml:"ticket"`
	Status      WorkflowType `yaml:"status"`
	Amount      types.String `yaml:"amount"`
	PaymentType types.String `yaml:"type"`
}

type DeleteTicket struct {
	Ticket types.Int `yaml:"ticket"`
}

func Processor(action string, order string, query *database.IQuery) error {
	var err error
	switch action {
	case "new_payment_flow":
		err = NewPaymentFlow(order, query)
	case "add_payment":
		err = DoAddPayment(order, query)
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

func UnmarshallAddPayment(order []byte) (AddPaymentStatus, error) {
	var change AddPaymentStatus
	var err error
	err = yaml.Unmarshal(order, &change)
	if err != nil {
		return change, err
	}
	return change, nil
}

func NewPaymentFlow(order string, query *database.IQuery) error {
	handler, err := UnmarshallTicket([]byte(order))
	if err != nil {
		return err
	}

	start := types.NewString(constants.Start)
	images := types.NewStringArray([]string{})
	roles := types.NewStringArray([]string{})
	taStory, err := models.CreateTaStory(*query,
		handler.Name,
		start,
		handler.Summary,
	)
	if err != nil {
		return err
	}
	models.CreateTaTracking(*query,
		taStory.PkStory,
		start,
		images,
		roles,
	)
	return nil
}

func DoAddPayment(order string, query *database.IQuery) error {
	handler, err := UnmarshallAddPayment([]byte(order))
	if err != nil {
		return err
	}

	taStory, err := models.GetTaStoryByPkStory(*query,
		handler.Ticket)
	if err != nil {
		return err
	}

	images := types.NewStringArray([]string{})
	roles := types.NewStringArray([]string{handler.Amount.String(), handler.PaymentType.String()})
	status := types.NewString(handler.Status.String())
	taStory.Update(*query,
		types.NilString(),
		status,
		types.NilString())

	models.CreateTaTracking(*query,
		taStory.PkStory,
		status,
		images,
		roles,
	)

	return nil
}

func DoDeleteTicket(order string, query *database.IQuery) error {
	handler, err := UnmarshallAddPayment([]byte(order))
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
