package service

import (
	"context"
	projpb "store/api/proj"
)

func (t *WorkFlowService) ReplaceWorkflow(ctx context.Context, request *projpb.Workflow) (*projpb.Workflow, error) {

	_, err := t.data.Mongo.Workflow.ReplaceByID(ctx, request.XId, request)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
