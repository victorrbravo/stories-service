package queries

import (
	"context"
	db "github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/pkg/workflow/common"
)

type Query[T db.Scanner] interface {
	Count(ctx context.Context, wf *common.Workflow[T], id string) error // id field is the name of task
	Percentage(ctx context.Context, wf *common.Workflow[T], id string) error
	TimeField(ctx context.Context, wf *common.Workflow[T], id string) error
	RoleField(ctx context.Context, wf *common.Workflow[T], id string) error
	Generate(ctx context.Context, wf *common.Workflow[T]) error
	StatusesList(ctx context.Context, wf *common.Workflow[T]) string
	Error() error
}
