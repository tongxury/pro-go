package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
)

func (t ProjAdminService) GetTemplate(ctx context.Context, request *projpb.GetGetTemplateRequest) (*projpb.Resource, error) {

	r, err := t.data.Mongo.Template.FindByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	list, err := t.data.Mongo.TemplateSegment.List(ctx, mgz.Filter().EQ("root._id", r.XId).B())
	if err != nil {
		return nil, err
	}

	r.Segments = list

	return r, nil
}
