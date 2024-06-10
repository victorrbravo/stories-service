package constants

const (
	ToDo                    string = "ToDo"
	Progress                string = "Progress"
	Done                    string = "Done"
	Start                   string = "Start"
	True                    string = "true"
	False                   string = "false"
	LimitPreffix                   = 3
	FlowTipsSuffix                 = "tips"
	SubSuffix                      = "sub"
	SummaryDefault                 = "No summary"
	EmptyStepsArray                = "[]"
	SetTipsIDNullInTracking        = "UPDATE ta_tracking SET fk_story_tips = NULL WHERE pk_tracking = $1"
	// actions
	AddFlowStepAction = "add_flow_step"
)
