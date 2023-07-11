package query

import (
	"context"
	"stories/pkg/workflow/common"
)

type Query interface {
	Count(ctx context.Context, wf *common.Workflow, id string) error // id field is the name of task
	Percentage(ctx context.Context, wf *common.Workflow, id string) error
	TimeField(ctx context.Context, wf *common.Workflow, id string) error
	RoleField(ctx context.Context, wf *common.Workflow, id string) error
	Generate(ctx context.Context, wf *common.Workflow) error
	StatusesList(ctx context.Context, wf *common.Workflow) string
}
