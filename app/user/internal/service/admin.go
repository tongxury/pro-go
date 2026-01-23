package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "store/api/admin"
	typepb "store/api/admin/types"
	"store/app/user/internal/biz"
	"store/pkg/krathelper"
)

type AdminService struct {
	pb.UnimplementedAdminServer
	adminBiz *biz.AdminBiz
}

func NewAdminService(adminBiz *biz.AdminBiz) *AdminService {
	return &AdminService{
		adminBiz: adminBiz,
	}
}

func (s *AdminService) GetNotifications(ctx context.Context, req *emptypb.Empty) (*typepb.Notification, error) {

	userId := krathelper.FindUserId(ctx)

	item, err := s.adminBiz.GetNotification(ctx, userId)
	if err != nil {
		log.Errorw("GetNotification err", err)
	}

	return item, nil
}
