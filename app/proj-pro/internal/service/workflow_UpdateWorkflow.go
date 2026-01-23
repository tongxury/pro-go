package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
)

func (t *WorkFlowService) UpdateWorkflow(ctx context.Context, request *projpb.UpdateWorkflowRequest) (*projpb.Workflow, error) {

	switch request.Action {
	case "updateSettings":

		_, err := t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, request.Id,
			mgz.Op().
				Set("dataBus.settings", request.Data.Settings),
		)
		if err != nil {
			return nil, err
		}
	case "cancel":
		err := t.workflow.Cancel(ctx, []string{request.Id})
		if err != nil {
			return nil, err
		}
	case "resume":

		err := t.workflow.Resume(ctx, []string{request.Id})
		if err != nil {
			return nil, err
		}
	case "updateField":
		op := mgz.Op()

		for k, v := range request.Kv {
			op = op.Set(k, v)
		}

		t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, request.Id,
			op)

	}

	return nil, nil
}
