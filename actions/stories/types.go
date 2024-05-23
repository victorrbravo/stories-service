package stories

import "stories/types"

type OrderInsert string

const (
	DateTime OrderInsert = "DATETIME"
	None     OrderInsert = "NONE"
)

type GetStoriesAction struct {
	Email types.String `yaml:"email"`
}

type PresignAction struct {
	Name types.String `yaml:"name"`
}

type SuggestAction struct {
	ImageID types.String `yaml:"image_id"`
}

type DeleteFlowStepAction struct {
	ID         types.Int      `yaml:"id"`
	DeleteType DeleteStepType `yaml:"delete_type"`
}

type Ticket struct {
	Name    types.String `yaml:"name"`
	Summary types.String `yaml:"summary"`
	UserID  types.Int    `yaml:"user_id"`
	Email   types.String `yaml:"email"`
}

type ShowStepTips struct {
	StepID types.Int `yaml:"step_id"`
}
type ChangeNameOfTask struct {
	StepID types.Int    `yaml:"step_id"`
	Name   types.String `yaml:"name"`
}
type ChangeStepField struct {
	StepID    types.Int         `yaml:"step_id"`
	Amount    types.String      `yaml:"amount_spent"`
	UpdatedAt types.String      `yaml:"updated_at"`
	Position  types.Coordinates `yaml:"position"`
}

type Tip struct {
	ID  types.Int    `yaml:"id"`
	Tip types.String `yaml:"tip"`
}

type AddStoryStatus struct {
	Ticket    types.Int         `yaml:"ticket"`
	Status    types.String      `yaml:"status"`
	Amount    types.String      `yaml:"amount"`
	Images    types.StringArray `yaml:"images"`
	Videos    types.StringArray `yaml:"videos"`
	StepType  types.String      `yaml:"type"`
	UserID    types.Int         `yaml:"user_id"`
	Summary   types.String      `yaml:"summary"`
	MediaDate types.String      `yaml:"media_date"`
	Position  types.Coordinates `yaml:"position"`
}
type StepTips struct {
	StepID types.Int         `yaml:"step_id"`
	Tips   types.StringArray `yaml:"tips"`
}

type DeleteTicket struct {
	Ticket types.Int `yaml:"ticket"`
}

type DeleteTipsByStep struct {
	StepID types.Int `yaml:"step_id"`
	Number types.Int `yaml:"number"`
}

type DeleteStepType string

const (
	TipsStepType DeleteStepType = "TIPS"
	SubStepType  DeleteStepType = "SUB"
)
