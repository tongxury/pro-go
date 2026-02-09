package service

import (
	"context"
	projpb "store/api/proj"
)

func (t ProjService) XGetResourceSegment(ctx context.Context, request *projpb.XGetResourceSegmentRequest) (*projpb.ResourceSegment, error) {

	id, err := t.data.Mongo.TemplateSegment.FindByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return id, nil
}
