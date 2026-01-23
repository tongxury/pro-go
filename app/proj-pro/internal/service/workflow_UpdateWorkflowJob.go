package service

import (
	"context"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/biz"
	"store/pkg/clients/mgz"
)

func (t *WorkFlowService) UpdateWorkflowJob(ctx context.Context, req *projpb.UpdateWorkflowJobRequest) (*projpb.Workflow, error) {

	switch req.Action {
	case "retry":

		//dataKeyStr := req.Params["dataKeys"]
		//dataKeys := strings.Split(dataKeyStr, ",")

		err := t.workflow.Retry(ctx, req.Id, req.Index)
		if err != nil {
			return nil, err
		}
	case "confirm":

		err := t.workflow.Confirm(ctx, req.Id, req.Index, req.RunImmediately)
		if err != nil {
			return nil, err
		}
	case "back":

		err := t.workflow.Back(ctx, req.Id, req.Index)
		if err != nil {
			return nil, err
		}
	case "updateData":
		return t.updateWorkflowJobData(ctx, req)
	}

	return nil, nil
}

func (t *WorkFlowService) updateWorkflowJobData(ctx context.Context, req *projpb.UpdateWorkflowJobRequest) (*projpb.Workflow, error) {

	switch req.GetDataKey() {
	case "keyFrames":
		_, err := t.data.Mongo.Workflow.UpdateByIDIfExists(ctx,
			req.Id,
			mgz.Op().
				SetListItem("jobs", int(req.Index), "dataBus.keyFrames", req.Data.KeyFrames).
				SetListItem("jobs", int(req.Index), "status", biz.JobStatusRunning),
		)
		if err != nil {
			return nil, err
		}
	case "videoGenerations":
		_, err := t.data.Mongo.Workflow.UpdateByIDIfExists(ctx,
			req.Id,
			mgz.Op().
				SetListItem("jobs", int(req.Index), "dataBus.videoGenerations", req.Data.VideoGenerations).
				SetListItem("jobs", int(req.Index), "status", biz.JobStatusRunning),
		)
		if err != nil {
			return nil, err
		}

	case "segmentScript":
		_, err := t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, req.Id,
			mgz.Op().
				SetListItem("jobs", int(req.Index), "dataBus.segmentScript", req.Data.SegmentScript),
		)
		if err != nil {
			return nil, err
		}

	case "remix":
		_, err := t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, req.Id,
			mgz.Op().
				SetListItem("jobs", int(req.Index), "dataBus.remix", req.Data.Remix).
				SetListItem("jobs", int(req.Index), "status", biz.JobStatusRunning),
		)
		if err != nil {
			return nil, err
		}

	}

	return nil, nil
}
