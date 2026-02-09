package service

import (
	"context"
	projpb "store/api/proj"
	ucpb "store/api/usercenter"
	"store/pkg/sdk/third/bytedance/vikingdb"

	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjService) XUpdateResourceSegment(ctx context.Context, req *projpb.XUpdateResourceSegmentRequest) (*projpb.ResourceSegment, error) {

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
		sg.Root.User = &ucpb.User{XId: req.UserId}

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
