package service

import (
	"context"
	creditpb "store/api/credit"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"store/pkg/sdk/helper"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (t ProjService) ReCreateAsset(ctx context.Context, req *projpb.ReCreateAssetRequest) (*projpb.Asset, error) {
	userId := krathelper.RequireUserId(ctx)

	state, err := t.data.GrpcClients.CreditClient.XGetCreditState(ctx, &creditpb.XGetCreditStateRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	if state.Balance < 10 {
		return nil, errors.BadRequest("exceeded", "积分不足，请联系客服充值")
	}

	log.Debugw("XGetCreditState state", state)

	refAsset, err := t.data.Mongo.Asset.FindByID(ctx, req.AssetId)
	if err != nil {
		return nil, err
	}

	videoId := helper.FindInStringMap(refAsset.Extra.GetContext(), "videoId")
	if videoId == "" {
		return nil, errors.BadRequest("videoIdEmpty", "videoId is empty")
	}

	groupId := refAsset.Group.GetId()
	if groupId == "" {

		groupId = primitive.NewObjectID().Hex()

		_, err = t.data.Mongo.Asset.UpdateByIDIfExists(ctx,
			req.AssetId,
			mgz.Op().Set("group.id", groupId),
		)
		if err != nil {
			return nil, err
		}
	}

	refAsset.Group.RefAssetId = req.AssetId

	newAsset, err := t.data.Mongo.Asset.Insert(ctx, &projpb.Asset{
		Commodity: refAsset.Commodity,
		Segment:   refAsset.Segment,
		UserId:    userId,
		Status:    "generating",
		CreatedAt: time.Now().Unix(),
		Prompt:    refAsset.Prompt,
		//Extra: &projpb.Asset_Extra{
		//	Context: refAsset.Extra.Context,
		//},
		PromptAddition: req.Prompt,
		Group: &projpb.Asset_Group{
			Id:         groupId,
			RefAssetId: req.AssetId,
		},
	})

	if err != nil {
		return nil, err
	}

	return newAsset, nil
}
func (t ProjService) CreateAsset(ctx context.Context, req *projpb.CreateAssetRequest) (*projpb.Asset, error) {

	userId := krathelper.RequireUserId(ctx)

	state, err := t.data.GrpcClients.CreditClient.XGetCreditState(ctx, &creditpb.XGetCreditStateRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	if state.Balance < 10 {
		return nil, errors.BadRequest("exceeded", "积分不足，请联系客服充值")
	}

	log.Debugw("XGetCreditState state", state)

	segment := req.Segment // 模板库的segment是没有 id的
	if segment == nil {
		sg, err := t.data.GrpcClients.ProjAdminClient.GetResourceSegment(ctx, &projpb.GetResourceSegmentRequest{
			Id: req.TemplateId,
		})
		if err != nil {
			return nil, err
		}

		segment = sg
	}

	if segment == nil {
		return nil, errors.BadRequest("not found template", "")
	}

	groupId := primitive.NewObjectID().Hex()
	if req.BaseAssetId != "" {
		baseAsset, err := t.data.Mongo.Asset.GetById(ctx, req.BaseAssetId)
		if err != nil {
			return nil, err
		}

		videoId := helper.FindInStringMap(baseAsset.Extra.GetContext(), "videoId")
		if videoId == "" {
			return nil, errors.BadRequest("videoIdEmpty", "videoId is empty")
		}

		groupId = baseAsset.Group.GetId()
	}

	// 重新拉一次 commodity 获取最新数据 或防止数据缺失
	commodity, err := t.data.Mongo.Commodity.GetById(ctx, req.Commodity.GetID())
	if err != nil {
		return nil, err
	}

	newAsset, err := t.data.Mongo.Asset.Insert(ctx, &projpb.Asset{
		Commodity:      commodity,
		Segment:        segment,
		UserId:         userId,
		Status:         "created",
		CreatedAt:      time.Now().Unix(),
		PromptAddition: req.PromptAddition,
		Prompts:        req.Prompts,
		Group: &projpb.Asset_Group{
			Id:          groupId,
			BaseAssetId: req.BaseAssetId,
		},
	})
	if err != nil {
		return nil, err
	}

	return newAsset, nil

}
