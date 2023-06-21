package query

import (
	"context"
	"stories/pkg/workflow/common"
)

type Query interface {
	Count(wf *common.Workflow, index int) error
	TimeField(ctx context.Context, wf *common.Workflow, id string) error
	RoleField(ctx context.Context, wf *common.Workflow, id string) error
	Generate(ctx context.Context, wf *common.Workflow) error
	StatusesList(ctx context.Context, wf *common.Workflow) string
}
