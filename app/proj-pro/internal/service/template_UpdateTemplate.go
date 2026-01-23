package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/third/bytedance/vikingdb"

	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjService) UpdateTemplate(ctx context.Context, params *projpb.UpdateTemplateRequest) (*projpb.Resource, error) {

	//userId := krathelper.RequireUserId(ctx)

	var err error
	switch params.Action {
	case "delete":
		//err = t.data.Elastics.Delete(ctx, projpb.ESIndexItems, params.Id)

		err := t.data.VikingDB.Delete(ctx, vikingdb.DeleteRequest{
			CollectionName: "template_commodity_coll",
			IDs:            []string{params.Id},
		})
		if err != nil {
			return nil, err
		}

		err = t.data.Mongo.Template.DeleteByID(ctx, params.Id)
	case "refresh":

		_, err := t.data.Mongo.TemplateSegment.Delete(ctx, bson.M{"root.id": params.Id})
		if err != nil {
			return nil, err
		}

		_, err = t.data.Mongo.Template.UpdateByIDXX(ctx, params.Id, mgz.Set(bson.M{"status": "created", "segments": nil}))

	}

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t ProjService) XUpdateTemplate(ctx context.Context, params *projpb.UpdateTemplateRequest) (*projpb.Resource, error) {

	var err error
	switch params.Action {
	case "delete":
		//err = t.data.Elastics.Delete(ctx, projpb.ESIndexItems, params.Id)
		err = t.data.Mongo.Template.DeleteByID(ctx, params.Id)
	case "refresh":

		_, err := t.data.Mongo.TemplateSegment.Delete(ctx, bson.M{"root.id": params.Id})
		if err != nil {
			return nil, err
		}

		_, err = t.data.Mongo.Template.UpdateByIDXX(ctx, params.Id, mgz.Set(bson.M{"status": "created", "segments": nil}))

	}

	if err != nil {
		return nil, err
	}

	return nil, nil
}
