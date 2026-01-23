package service

import (
	"context"
	"store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/third/bytedance/tos"
)

func (t ProjService) ListQualityAssets(ctx context.Context, params *projpb.ListQualityAssetsRequest) (*projpb.AssetList, error) {

	//userId := krathelper.RequireUserId(ctx)

	filter := mgz.Filter().
		EQ("userId", "690d6a669e5c05462c0e4165").
		EQ("favorite", true)

	list, c, err := t.data.Mongo.Asset.ListAndCount(ctx,
		filter.B(),
		mgz.Find().
			Paging(params.Page, params.Size).
			SetFields(params.ReturnFields).
			SetSort("createdAt", -1).
			B(),
	)
	if err != nil {
		return nil, err
	}
	for _, x := range list {
		x.Extra = nil
		if x.CoverUrl == "" {
			x.CoverUrl = tos.Change(x.Commodity.GetMedias()[0].GetUrl())
		}
	}
	return &projpb.AssetList{
		List:    list,
		Total:   c,
		Page:    params.Page,
		Size:    params.Size,
		HasMore: params.Size*params.Page < c,
	}, err
}
