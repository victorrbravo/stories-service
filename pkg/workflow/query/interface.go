package query

import (
	"context"
	"stories/pkg/workflow/common"
)

type Query interface {
	Count(wf *common.Workflow, index int) error
	TimeField(wf *common.Workflow, index int) error
	Generate(ctx context.Context, wf *common.Workflow) error
}
