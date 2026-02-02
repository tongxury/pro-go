package biz

import (
	"context"
	projpb "store/api/proj"
)

type ExecuteResult struct {
	Status      string
	Cost        int
	SkipConfirm bool
	Error       string
}

type Options struct {
	JobState       *projpb.Job
	WorkflowState  *projpb.Workflow
	RunImmediately bool
}

type IJob interface {
	GetName() string
	Initialize(ctx context.Context, options Options) error
	Execute(ctx context.Context, jobState *projpb.Job, workflowState *projpb.Workflow) (status *ExecuteResult, err error)
}
