package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
)

func (t ProjService) GetTemplate(ctx context.Context, request *projpb.GetTemplateRequest) (*projpb.Resource, error) {
	return t.XGetTemplate(ctx, request)
}

func (t ProjService) XGetTemplate(ctx context.Context, request *projpb.GetTemplateRequest) (*projpb.Resource, error) {
	root, err := t.data.Mongo.Template.FindByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	list, err := t.data.Mongo.TemplateSegment.List(ctx,
		mgz.Filter().EQ("root._id", root.XId).B(),
		mgz.Find().SetSort("timeStart", 1).B(),
	)
	if err != nil {
		return nil, err
	}

	root.Segments = list

	return root, nil
}
