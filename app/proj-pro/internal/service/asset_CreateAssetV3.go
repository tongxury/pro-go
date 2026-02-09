package service

import (
	"context"
	creditpb "store/api/credit"
	projpb "store/api/proj"
	usercenter "store/api/usercenter"
	"store/pkg/krathelper"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

func (t ProjService) CreateAssetV3(ctx context.Context, req *projpb.CreateAssetV3Request) (*projpb.Asset, error) {

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

	//
	workflow, err := t.workflow.CreateVideoGenerationWorkflow(ctx,
		userId,
		&projpb.SegmentScript{
			Inspiration: req.Prompt,
			Images:      req.Images,
		})
	if err != nil {
		return nil, err
	}

	newAsset, err := t.data.Mongo.Asset.Insert(ctx, &projpb.Asset{
		User:      &usercenter.User{XId: userId},
		Status:    "created",
		CreatedAt: time.Now().Unix(),
		Category:  "videoGeneration",
		Workflow: &projpb.Workflow{
			XId:    workflow.XId,
			UserId: userId,
		},
	})
	if err != nil {
		return nil, err
	}

	return newAsset, nil

}
