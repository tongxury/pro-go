package service

import (
	"context"
	projpb "store/api/proj"
)

func (t *WorkFlowService) CreateWorkflow(ctx context.Context, request *projpb.CreateWorkflowRequest) (*projpb.Workflow, error) {

	//userId := krathelper.RequireUserId(ctx)
	userId := "693b91f82b271bf02ac1e624"

	workflow, err := t.workflow.CreateVideoReplication3Workflow(ctx, userId, &projpb.DataBus{
		Commodity: request.Commodity,
	})
	if err != nil {
		return nil, err
	}

	return workflow, nil
}
