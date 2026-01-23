package service

import (
	"context"
	projpb "store/api/proj"
)

func (t ProjAdminService) GetResourceSegment(ctx context.Context, request *projpb.GetResourceSegmentRequest) (*projpb.ResourceSegment, error) {

	id, err := t.data.Mongo.TemplateSegment.FindByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return id, nil
}
