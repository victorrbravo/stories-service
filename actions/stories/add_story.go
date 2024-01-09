package stories

import (
	"context"
	"gopkg.in/yaml.v3"
	"stories/actions/constants"
	"stories/database"
	"stories/generated/models"
	"stories/pkg/utils"
	"stories/types"
)

type Ticket struct {
	Name    types.String `yaml:"name"`
	Summary types.String `yaml:"summary"`
}

type AddStoryStatus struct {
	Ticket   types.Int    `yaml:"ticket"`
	Status   types.String `yaml:"status"`
	Amount   types.String `yaml:"amount"`
	Image1   types.String `yaml:"image"`
	StepType types.String `yaml:"type"`
}

type DeleteTicket struct {
	Ticket types.Int `yaml:"ticket"`
}

func FlowProcessor(action string, order string, query *database.IQuery) error {
	var err error
	ctx := context.Background()
	switch action {
	case "new_story_flow":
		err = NewStoryFlow(order, query)
	case "add_step":
		_, err = DoAddStep(order, query)
	case "add_flow_step":
		err = DoAddStepFlowAtEnd(ctx, order, query)
	case "delete_ticket":
		err = DoDeleteStep(order, query)
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

func UnmarshallAddStep(order []byte, query *database.IQuery) (AddStoryStatus, error) {
	var change AddStoryStatus
	var err error
	err = yaml.Unmarshal(order, &change)
	if err != nil {
		return change, err
	}
	return change, nil
}

func NewStoryFlow(order string, query *database.IQuery) error {
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

func DoAddStepFlowAtEnd(ctx context.Context, order string, query *database.IQuery) error {
	addStep, err := UnmarshallAddStep([]byte(order), query)
	if err != nil {
		return err
	}
	var addTask AddTask
	wf, err := GetWorkflow(1, query)
	if err != nil {
		return err
	}
	addTask.Name = types.NewString(addStep.Status.String())
	lastTask, err := wf.Last(ctx)
	if err != nil {
		return err
	}
	addTask.PkFlow = 1
	addTask.Source = types.NewString(lastTask.Name)
	addTask.Sentence = types.NewString(wf.Sentence(addTask.Name.String()))
	utils.SetDefaults(&addTask)
	wf, err = DoAddTaskIntl(addTask, query)
	if err != nil {
		return err
	}
	err = DoAddStepIntl(addStep, query)
	if err != nil {
		return err
	}
	return nil
}

func DoAddStepFlow(order string, query *database.IQuery) error {
	handler, err := UnmarshallAddStep([]byte(order), query)
	if err != nil {
		return err
	}

	taStory, err := models.GetTaStoryByPkStory(*query,
		handler.Ticket)
	if err != nil {
		return err
	}

	images := types.NewStringArray([]string{handler.Image1.String()})
	roles := types.NewStringArray([]string{handler.Amount.String(), handler.StepType.String()})
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

func DoAddStepIntl(status AddStoryStatus, query *database.IQuery) error {
	taStory, err := models.GetTaStoryByPkStory(*query,
		status.Ticket)
	if err != nil {
		return err
	}

	images := types.NewStringArray([]string{status.Image1.String()})
	roles := types.NewStringArray([]string{status.Amount.String(), status.StepType.String()})
	statusStr := types.NewString(status.Status.String())
	taStory.Update(*query,
		types.NilString(),
		statusStr,
		types.NilString())

	models.CreateTaTracking(*query,
		taStory.PkStory,
		statusStr,
		images,
		roles,
	)

	return nil
}

func DoAddStep(order string, query *database.IQuery) (AddStoryStatus, error) {
	handler, err := UnmarshallAddStep([]byte(order), query)
	if err != nil {
		return handler, err
	}
	err = DoAddStepIntl(handler, query)
	if err != nil {
		return handler, err
	}
	return handler, nil
}

func DoDeleteStep(order string, query *database.IQuery) error {
	handler, err := UnmarshallAddStep([]byte(order), query)
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
