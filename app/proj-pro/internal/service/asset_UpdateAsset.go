package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (t ProjService) UpdateAsset(ctx context.Context, request *projpb.UpdateAssetRequest) (*projpb.Asset, error) {

	asset, err := t.data.Mongo.Asset.FindByID(ctx, request.Id)
	if err != nil {
		return nil, errors.BadRequest("invalidAssetId", "")
	}

	switch request.Action {

	case "favorite":
		_, err = t.data.Mongo.Asset.C().UpdateOne(ctx,
			bson.M{"_id": mgz.ObjectId(asset.XId)},
			mongo.Pipeline{
				{{"$set", bson.D{{"favorite", bson.D{{"$not", "$favorite"}}}}}},
				{{"$set", bson.D{{"attrs.favorite", bson.D{{"$not", "$attrs.favorite"}}}}}},
			},
		)
		if err != nil {
			return nil, err
		}
	case "attitude":
		_, err = t.data.Mongo.Asset.UpdateByIDIfExists(ctx,
			asset.XId,
			mgz.Op().Set("attrs.attitude", request.Params["attitude"]),
		)
	case "remark":
		_, err = t.data.Mongo.Asset.UpdateByIDIfExists(ctx,
			asset.XId,
			mgz.Op().Set("remark", request.Params["remark"]),
		)
	case "delete":
		err = t.data.Mongo.Asset.DeleteByID(ctx,
			asset.XId,
			mgz.DeleteOptions{
				Reserve: true,
			},
		)
	case "generatePrompt":
		_, err = t.data.Mongo.Asset.UpdateByIDIfExists(ctx, asset.XId,
			mgz.Op().
				Sets(bson.M{
					"status": "promptGenerating",
				}),
		)
		if err != nil {
			return nil, err
		}
	case "generate":

		_, err = t.data.Mongo.Asset.UpdateByIDIfExists(ctx,
			asset.XId,
			mgz.Op().
				Sets(bson.M{
					"status": "generating",
				}),
		)
		if err != nil {
			return nil, err
		}
	case "regenerate":
		_, err = t.data.Mongo.Asset.UpdateByIDIfExists(ctx, asset.XId,
			mgz.Op().
				Sets(bson.M{
					"status":         "generating",
					"promptAddition": request.Text,
				}),
		)
		if err != nil {
			return nil, err
		}
	case "updatePrompt":

		if strings.TrimSpace(request.Text) == "" {
			return nil, errors.BadRequest("cannotSetEmptyPrompt", "")
		}

		_, err = t.data.Mongo.Asset.UpdateByIDIfExists(ctx, asset.XId,
			mgz.Op().
				Sets(bson.M{
					"prompt": request.Text,
				}),
		)
		if err != nil {
			return nil, err
		}
	}

	return &projpb.Asset{}, nil
}
