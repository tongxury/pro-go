package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/krathelper"
)

func (t ProjService) GetResourceSegment(ctx context.Context, request *projpb.GetResourceSegmentRequest) (*projpb.ResourceSegment, error) {

	id, err := t.data.Mongo.TemplateSegment.FindByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	userId, _ := krathelper.GetUserId(ctx)
	if userId != "" {
		t.fillCollectionStatus(ctx, userId, []*projpb.ResourceSegment{id})
	}

	return id, nil
}
