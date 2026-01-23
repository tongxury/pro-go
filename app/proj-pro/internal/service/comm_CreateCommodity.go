package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/krathelper"
	"store/pkg/sdk/helper"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjService) CreateCommodity(ctx context.Context, params *projpb.CreateCommodityRequest) (*projpb.Commodity, error) {

	userId := krathelper.RequireUserId(ctx)

	// 获取商品信息
	metadata, err := t.data.Douyin.GetCommodityMetadataByAPI(ctx, params.Url)
	if err != nil {
		log.Errorw("GetCommodityMetadata err", err, "params", params)
		return nil, errors.New(0, "invalidUrl", "无效链接，请检查并重试")
	}

	newCommodity := &projpb.Commodity{
		Url:       params.Url,
		Medias:    params.Medias,
		Status:    "created",
		UserId:    userId,
		CreatedAt: time.Now().Unix(),
		Title:     metadata.Title,
		Brand:     metadata.Brand,
		Images:    helper.SubSlice(metadata.Images, 20),
	}

	_, _, err = t.data.Mongo.Commodity.InsertNX(ctx,
		newCommodity,
		bson.M{"userId": userId,
			"url": params.Url,
		},
	)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
