package service

import (
	"context"
	"fmt"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
)

func (t *WorkFlowService) CreateRemixTask(ctx context.Context, req *projpb.CreateRemixTaskRequest) (*projpb.RemixTask, error) {

	if len(req.Items) == 0 {
		return nil, fmt.Errorf("items is empty")
	}

	workflow, err := t.data.Mongo.Workflow.GetById(ctx, req.WorkflowId)
	if err != nil {
		return nil, err
	}

	if workflow.Name != "VideoReplication3" {
		return nil, fmt.Errorf("workflow name is not VideoReplication3")
	}

	var vgs []*projpb.VideoGeneration
	for _, item := range req.Items {
		vgs = append(vgs, &projpb.VideoGeneration{
			Url:      item.Url,
			Subtitle: item.Subtitle,
			Status:   "completed",
		})
	}

	_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, req.WorkflowId,
		mgz.Op().
			Set("jobs.3.dataBus.videoGenerations", vgs).
			Set("jobs.3.status", "completed").
			Set("jobs.4.status", "running").
			Set("jobs.4.dataBus.remix", nil).
			Set("status", "running"),
	)
	if err != nil {
		return nil, err
	}

	return &projpb.RemixTask{
		TaskId: req.WorkflowId,
		Status: "running", // Initial status
	}, nil
}

func (t *WorkFlowService) GetRemixTask(ctx context.Context, req *projpb.GetRemixTaskRequest) (*projpb.RemixTask, error) {

	workflow, err := t.data.Mongo.Workflow.GetById(ctx, req.TaskId)
	if err != nil {
		return nil, err
	}

	var url string
	for _, job := range workflow.Jobs {
		if job.Name == "VideoSegmentsRemixJob" {
			if job.Status == "completed" {
				url = job.DataBus.Remix.Url
				return &projpb.RemixTask{
					TaskId: req.TaskId,
					Status: "completed",
					Url:    url,
				}, nil
			} else {
				return &projpb.RemixTask{
					TaskId: req.TaskId,
					Status: "running",
				}, nil
			}
		}
	}

	return &projpb.RemixTask{
		TaskId: req.TaskId,
		Status: "running",
	}, nil
}
