package service

import (
	"context"
	projpb "store/api/proj"
	ucpb "store/api/usercenter"
	"store/pkg/krathelper"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (t ProjService) UpdateResourceSegment(ctx context.Context, req *projpb.UpdateResourceSegmentRequest) (*emptypb.Empty, error) {
	userId := krathelper.RequireUserId(ctx)
	if req.Id == "" {
		return nil, errors.BadRequest("missing_id", "id is required")
	}

	if req.Action == "collect" {
		// Verify segment exists
		_, err := t.data.Mongo.TemplateSegment.FindByID(ctx, req.Id)
		if err != nil {
			return nil, errors.NotFound("segment_not_found", "resource segment not found")
		}

		_, _, err = t.data.Mongo.ResourceSegmentCollection.InsertNX(ctx, &projpb.ResourceSegmentCollectionItem{
			User:            &ucpb.User{XId: userId},
			ResourceSegment: &projpb.ResourceSegment{XId: req.Id},
			CreatedAt:       time.Now().Unix(),
		}, bson.M{
			"user._id":            userId,
			"resourceSegment._id": req.Id,
		})
		return &emptypb.Empty{}, err
	}

	if req.Action == "cancel" {
		_, err := t.data.Mongo.ResourceSegmentCollection.Delete(ctx, bson.M{
			"user._id": userId,
			"$or": []bson.M{
				{"resourceSegment._id": req.Id},
				{"resourceSemgent._id": req.Id},
			},
		})
		return &emptypb.Empty{}, err
	}

	return nil, errors.BadRequest("invalid_action", "invalid action")
}
