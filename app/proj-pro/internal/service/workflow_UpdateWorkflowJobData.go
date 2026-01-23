package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
)

func (t *WorkFlowService) UpdateWorkflowJobData(ctx context.Context, req *projpb.UpdateWorkflowJobDataRequest) (*projpb.Workflow, error) {

	switch req.Name {
	case "keyFrames":
		// 修改数据可能要重新执行job  这里后面优化
		_, err := t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, req.Id,
			mgz.Op().
				Set("dataBus.keyFrames", req.Data.KeyFrames),
		)
		if err != nil {
			return nil, err
		}

	case "segmentScript":
		_, err := t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, req.Id,
			mgz.Op().Set("dataBus.segmentScript", req.Data.SegmentScript),
		)
		if err != nil {
			return nil, err
		}

	}

	return nil, nil
}
