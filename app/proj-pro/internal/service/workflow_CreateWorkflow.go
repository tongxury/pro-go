package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/krathelper"
)

func (t *WorkFlowService) CreateWorkflow(ctx context.Context, request *projpb.CreateWorkflowRequest) (*projpb.Workflow, error) {

	userId := krathelper.RequireUserId(ctx)

	workflow, err := t.workflow.CreateVideoReplication3Workflow(ctx, userId, &projpb.DataBus{
		Commodity: request.Commodity,
	})
	if err != nil {
		return nil, err
	}

	return workflow, nil
}
