package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/third/bytedance/tos"
)

func (t ProjService) GetAsset(ctx context.Context, request *projpb.GetAssetRequest) (*projpb.Asset, error) {
	return t.XGetAsset(ctx, &projpb.XGetAssetRequest{
		Id:           request.Id,
		ReturnFields: request.ReturnFields,
	})
}

func (t ProjService) XGetAsset(ctx context.Context, request *projpb.XGetAssetRequest) (*projpb.Asset, error) {

	asset, err := t.data.Mongo.Asset.FindByID(ctx, request.Id, mgz.Find().SetFields(request.GetReturnFields()).B())
	if err != nil {
		asset, err = t.data.Mongo.Asset.FindOne(ctx, mgz.Filter().EQ("workflow._id", request.Id).B(), mgz.Find().SetFields(request.GetReturnFields()).B())
		if err != nil {
			return nil, err
		}
	}

	if asset.GetWorkflow().GetID() != "" {
		asset.Workflow, err = t.data.Mongo.Workflow.GetById(ctx, asset.GetWorkflow().GetID())
		if err != nil {
			return nil, err
		}
	}

	list, err := t.data.Mongo.Feedback.List(ctx, mgz.Filter().EQ("targetId", asset.XId).B())
	if err != nil {
		return nil, err
	}

	asset.Url = tos.Change(asset.Url)
	asset.CoverUrl = tos.Change(asset.CoverUrl)

	asset.EnsureAttrs().FeedbackCount = int64(len(list))

	return asset, nil
}
