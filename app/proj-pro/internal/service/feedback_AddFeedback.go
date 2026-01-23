package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/krathelper"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (t ProjService) AddFeedback(ctx context.Context, request *projpb.AddFeedbackRequest) (*emptypb.Empty, error) {

	userId := krathelper.RequireUserId(ctx)

	_, err := t.data.Mongo.Feedback.Insert(ctx, &projpb.Feedback{
		Category:  request.Category,
		Url:       request.Url,
		Content:   request.Content,
		UserId:    userId,
		CreatedAt: time.Now().Unix(),
		Issues:    request.Issues,
		TargetId:  request.TargetId,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
