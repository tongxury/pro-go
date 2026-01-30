package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"strings"
)

func (t WorkFlowService) ListWorkflows(ctx context.Context, params *projpb.ListWorkflowsRequest) (*projpb.WorkflowList, error) {

	userId := krathelper.RequireUserId(ctx)

	filter := mgz.Filter()

	if userId != "" {
		filter = filter.EQ("userId", userId)
	}

	if params.Status != "" {
		filter = filter.InString("status", strings.Split(params.Status, ","))
	}

	list, c, err := t.data.Mongo.Workflow.ListAndCount(ctx,
		filter.B(),
		mgz.Find().
			Paging(params.Page, params.Size).
			SetSort("createdAt", -1).
			B(),
	)
	if err != nil {
		return nil, err
	}

	return &projpb.WorkflowList{
		List:    list,
		Total:   c,
		Page:    params.Page,
		Size:    params.Size,
		HasMore: params.Size*params.Page < c,
	}, nil

}
