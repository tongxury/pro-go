package service

import (
	"context"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/biz"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"store/pkg/sdk/helper"
)

func (t ProjService) ListOngoingIssues(ctx context.Context, req *projpb.ListOngoingIssuesRequest) (*projpb.IssueList, error) {

	userId := krathelper.RequireUserId(ctx)

	var issues []*projpb.Issue

	sessions, err := t.data.Mongo.Session.List(ctx,
		mgz.Filter().
			EQ("userId", userId).
			NotEQ("status", "completed").B(),

		mgz.Find().SetFields("_id,status,category,commodity.title").B(),
	)
	if err != nil {
		return nil, err
	}

	for _, x := range sessions {
		issues = append(issues, &projpb.Issue{
			XId:      x.XId,
			Category: helper.OrString(x.Category, "templateReplication"),
			Title:    x.Commodity.GetTitle(),
		})
	}

	workflows, err := t.data.Mongo.Workflow.List(ctx,
		mgz.Filter().
			EQ("userId", userId).
			And(
				mgz.Filter().NotEQ("status", biz.WorkflowStatusCompleted),
				mgz.Filter().NotEQ("status", biz.WorkflowStatusCanceled),
			).
			B(),
		mgz.Find().SetFields("_id,status,category,commodity.title").B(),
	)
	if err != nil {
		return nil, err
	}

	ids := mgz.Ids(workflows)
	if len(ids) > 0 {
		assets, err := t.data.Mongo.Asset.List(ctx, mgz.Filter().InString("workflow._id", ids).B())
		if err != nil {
			return nil, err
		}

		for _, x := range assets {
			issues = append(issues, &projpb.Issue{
				XId:      x.XId,
				Category: "segmentReplication",
				Title:    x.Commodity.GetTitle(),
			})
		}
	}

	return &projpb.IssueList{
		List: issues,
	}, nil
}
