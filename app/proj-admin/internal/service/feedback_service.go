package service

import (
	"context"
	projpb "store/api/proj"
)

func (t ProjAdminService) ListFeedbacks_(ctx context.Context, request *projpb.ListFeedbacksRequest_) (*projpb.FeedbackList, error) {

	templates, err := t.data.GrpcClients.ProjProClient.XListFeedbacks(ctx, &projpb.XListFeedbacksRequest{
		Page:     request.Page,
		Category: request.Category,
		Size:     request.Size,
		TargetId: request.TargetId,
	})
	if err != nil {
		return nil, err
	}
	return templates, nil
}
