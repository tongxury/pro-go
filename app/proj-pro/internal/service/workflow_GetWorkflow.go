package service

import (
	"context"
	projpb "store/api/proj"

	"github.com/go-kratos/kratos/v2/log"
)

func (t *WorkFlowService) GetWorkflow(ctx context.Context, request *projpb.GetWorkflowRequest) (*projpb.Workflow, error) {

	id, err := t.data.Mongo.Workflow.GetById(ctx, request.Id)
	if err != nil {
		log.Errorw("GetWorkflow err", err, "id", request.Id)
		return nil, err
	}

	return id, nil
}
