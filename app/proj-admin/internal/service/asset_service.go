package service

import (
	"context"
	projpb "store/api/proj"
	ucpb "store/api/usercenter"

	"github.com/go-kratos/kratos/v2/log"
)

func (t ProjAdminService) ListAssets_(ctx context.Context, request *projpb.ListAssetsRequest_) (*projpb.AssetList, error) {

	userId := request.UserId
	if request.UserPhone != "" {
		user, err := t.data.GrpcClients.UserCenterClient.XGetUserByPhone(ctx, &ucpb.XGetUserByPhoneRequest{
			Phone: request.UserPhone,
		})

		if err != nil {
			return nil, err
		}

		if user == nil {
			return nil, nil
		}

		log.Debugw("XGetUserByPhone ", "", "user", user)

		userId = user.XId
	}

	return t.data.GrpcClients.ProjProClient.XListAssets(ctx, &projpb.XListAssetsRequest{
		Page:         request.Page,
		Size:         request.Size,
		ReturnFields: request.ReturnFields,
		Status:       request.Status,
		UserId:       userId,
		Id:           request.Id,
	})
}

func (t ProjAdminService) GetAsset_(ctx context.Context, request *projpb.GetAssetRequest_) (*projpb.Asset, error) {
	return t.data.GrpcClients.ProjProClient.XGetAsset(ctx, &projpb.XGetAssetRequest{
		Id: request.Id,
	})
}
