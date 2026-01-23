package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/sdk/third/bytedance/vikingdb"

	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjAdminService) UpdateResourceSegment(ctx context.Context, req *projpb.UpdateResourceSegmentRequest) (*projpb.ResourceSegment, error) {

	var err error
	switch req.Action {
	case "delete":

		err := t.data.VikingDB.Delete(ctx, vikingdb.DeleteRequest{
			CollectionName: "segment_commodity_coll",
			IDs:            []string{req.Id},
		})
		if err != nil {
			return nil, err
		}

		err = t.data.Mongo.TemplateSegment.DeleteByID(ctx, req.Id)

		//_, err = t.data.Elastics.DeleteByRequest(ctx, projpb.ESIndexItemSegments, deletebyquery.Request{
		//	Query: elastics.NewTermQuery("root.id.keyword", req.RootId),
		//})
	case "refresh":

		sg, err := t.data.Mongo.TemplateSegment.FindByID(ctx, req.Id)
		if err != nil {
			return nil, err
		}

		if sg == nil {
			return nil, nil
		}

		_, err = t.data.Mongo.TemplateSegment.Delete(ctx, bson.M{"root.url": sg.Root.Url})

		sg.Status = "created"

		_, err = t.data.Mongo.TemplateSegment.InsertMany(ctx, sg)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return &projpb.ResourceSegment{}, nil
}
