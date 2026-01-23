package service

import (
	"context"
	"fmt"
	creditpb "store/api/credit"
	projpb "store/api/proj"
	ucpb "store/api/usercenter"
	"store/pkg/sdk/conv"
	"time"
)

func (t ProjAdminService) ListUsers_(ctx context.Context, request *projpb.ListUsersRequest_) (*ucpb.UserList, error) {
	return t.data.GrpcClients.UserCenterClient.XListUsers(ctx, &ucpb.XListUsersRequest{
		Page:    request.Page,
		Size:    request.Size,
		Keyword: request.Keyword,
	})
}

func (t ProjAdminService) GetUser_(ctx context.Context, request *projpb.GetUserRequest_) (*projpb.UserInfo, error) {
	user, err := t.data.GrpcClients.UserCenterClient.XGetUser(ctx, &ucpb.XGetUserRequest{
		Id: request.Id,
	})

	if err != nil {
		return nil, err
	}

	state, err := t.data.GrpcClients.CreditClient.XGetCreditState(ctx, &creditpb.XGetCreditStateRequest{
		UserId: request.Id,
	})
	if err != nil {
		return nil, err
	}

	return &projpb.UserInfo{
		User:        user,
		CreditState: state,
	}, nil
}

func (t ProjAdminService) UpdateUser_(ctx context.Context, request *projpb.UpdateUserRequest_) (*ucpb.User, error) {

	switch request.Action {
	case "recharge":
		_, err := t.data.GrpcClients.CreditClient.XRecharge(ctx, &creditpb.XRechargeRequest{
			UserId:   request.Id,
			Amount:   conv.Int64(request.Params["amount"]),
			Key:      fmt.Sprintf("%s_%d", request.Id, time.Now().Unix()),
			Category: "manual",
		})
		if err != nil {
			return nil, err
		}

	}

	return nil, nil
}
