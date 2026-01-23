package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"time"
)

func (t ProjService) CreateAssetV5(ctx context.Context, req *projpb.CreateAssetV5Request) (*projpb.Asset, error) {

	userId := krathelper.RequireUserId(ctx)

	newAsset := &projpb.Asset{
		UserId:    userId,
		Status:    "completed",
		CreatedAt: time.Now().Unix(),
		Category:  "custom",
		CoverUrl:  req.CoverUrl,
		Url:       req.Url,
	}

	id, _, err := t.data.Mongo.Asset.InsertNX(ctx, newAsset, mgz.Filter().EQ("url", req.Url).B())
	if err != nil {
		return nil, err
	}

	newAsset.XId = id
	return newAsset, nil

}
