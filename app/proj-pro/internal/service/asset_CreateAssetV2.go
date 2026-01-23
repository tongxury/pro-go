package service

import (
	"context"
	creditpb "store/api/credit"
	projpb "store/api/proj"
	"store/pkg/krathelper"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

func (t ProjService) CreateAssetV2(ctx context.Context, req *projpb.CreateAssetV2Request) (*projpb.Asset, error) {

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

	// 模版的segment 没有id  包括后面可能得其他模版
	segment := req.Segment
	if req.SegmentId != "" {
		segment, err = t.data.GrpcClients.ProjAdminClient.GetResourceSegment(ctx, &projpb.GetResourceSegmentRequest{
			Id: req.SegmentId,
		})
		if err != nil {
			return nil, err
		}
	}

	if segment == nil {
		return nil, errors.BadRequest("not found template", "")
	}
	//
	//if len(segment.Segments) == 0 {
	//	return nil, errors.BadRequest("not found segment", "请重新分析当前高光片段")
	//}

	// 重新拉一次 commodity 获取最新数据 或防止数据缺失
	commodity, err := t.data.Mongo.Commodity.GetById(ctx, req.CommodityId)
	if err != nil {
		return nil, err
	}

	//
	workflow, err := t.workflow.CreateVideoReplicationWorkflow(ctx, userId, req.WorkflowName, segment, commodity, req.InitialData)
	if err != nil {
		return nil, err
	}

	newAsset, err := t.data.Mongo.Asset.Insert(ctx, &projpb.Asset{
		Commodity: commodity,
		Segment:   segment,
		UserId:    userId,
		Status:    "created",
		Category:  "segmentReplication",
		CreatedAt: time.Now().Unix(),
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
