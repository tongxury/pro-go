package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
)

func (t ProjService) ListFeedbacks(ctx context.Context, request *projpb.ListFeedbacksRequest) (*projpb.FeedbackList, error) {

	userId := krathelper.RequireUserId(ctx)

	return t.XListFeedbacks(ctx, &projpb.XListFeedbacksRequest{
		Page:     request.Page,
		Category: request.Category,
		Size:     request.Size,
		TargetId: request.TargetId,
		UserId:   userId,
	})
}
func (t ProjService) XListFeedbacks(ctx context.Context, request *projpb.XListFeedbacksRequest) (*projpb.FeedbackList, error) {

	filter := mgz.Filter()

	if request.Category != "" {
		filter.EQ("category", request.Category)
	}
	if request.TargetId != "" {
		filter.EQ("targetId", request.TargetId)
	}

	if request.UserId != "" {
		filter.EQ("userId", request.UserId)
	}

	list, c, err := t.data.Mongo.Feedback.ListAndCount(ctx,
		filter.B(),
		mgz.Find().
			Paging(request.Page, request.Size).
			SetSort("_id", -1).
			B(),
	)
	if err != nil {
		return nil, err
	}

	return &projpb.FeedbackList{
		List:  list,
		Total: c,
	}, nil
}
