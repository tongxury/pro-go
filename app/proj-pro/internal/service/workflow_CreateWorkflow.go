package service

import (
	"context"
	projpb "store/api/proj"
)

func (t *WorkFlowService) CreateWorkflow(ctx context.Context, request *projpb.CreateWorkflowRequest) (*projpb.Workflow, error) {

	t.workflow.CreateVideoGenerationWorkflow()

	return nil, nil
}
