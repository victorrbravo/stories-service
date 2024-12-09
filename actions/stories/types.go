package stories

import (
	"stories/generated/models"
	"stories/types"
)

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

type MapResponse struct {
	Values map[string][]models.TaStory `json:"values"`
}

type Ticket struct {
	Name        types.String      `yaml:"name"`
	Summary     types.String      `yaml:"summary"`
	UserID      types.Int         `yaml:"user_id"`
	Email       types.String      `yaml:"email"`
	ContentType types.ContentType `yaml:"content_type"`
}

type ShowStepTips struct {
	StepID types.Int `yaml:"step_id"`
	TipsID types.Int `yaml:"tips_id"`
}

type ChangeNameOfTask struct {
	StepID types.Int    `yaml:"step_id" json:"step_id"`
	Name   types.String `yaml:"name" json:"name"`
}
type ChangeStepField struct {
	StepID    types.Int         `yaml:"step_id"`
	Amount    types.String      `yaml:"amount_spent"`
	UpdatedAt types.String      `yaml:"updated_at"`
	Position  types.Coordinates `yaml:"position"`
}

type MoveTaskInFlow struct {
	PkFlow  types.Int    `yaml:"id_flow"`
	PkStory types.Int    `yaml:"id_story"`
	Current types.String `yaml:"current"`
	Target  types.String `yaml:"target"`
}

type UpdateAttrInFlow struct {
	PkFlow  types.Int    `yaml:"id_flow"`
	PkStory types.Int    `yaml:"id_story"`
	Report  types.Bool   `yaml:"report"`
	Name    types.String `yaml:"name"`
}

type DataTip struct {
	Tips  []Tip `yaml:"tips" json:"tips"`
	Count int   `yaml:"count" json:"count"`
}

type Tip struct {
	ID  types.Int    `yaml:"id" json:"id"`
	Tip types.String `yaml:"tip" json:"tip"`
}

type AddStoryStatus struct {
	Ticket    types.Int         `yaml:"ticket,omitempty"`
	Status    types.String      `yaml:"status"`
	Amount    types.String      `yaml:"amount"`
	Images    types.StringArray `yaml:"images"`
	Videos    types.StringArray `yaml:"videos"`
	StepType  types.String      `yaml:"type"`
	UserID    types.Int         `yaml:"user_id"`
	Summary   types.String      `yaml:"summary"`
	MediaDate types.String      `yaml:"media_date"`
	Position  types.Coordinates `yaml:"position"`
	metaOrder types.String      `yaml:"meta_order"`
	Tags      types.StringArray `yaml:"tags"`
}

type StepTips struct {
	StepID types.Int         `yaml:"step_id"`
	Tips   types.StringArray `yaml:"tips"`
}

type StepSub struct {
	StepID types.Int      `yaml:"step_id"`
	Task   AddStoryStatus `yaml:"task"`
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
