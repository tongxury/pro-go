package service

import (
	"context"
	projpb "store/api/proj"
)

func (t ProjAdminService) GetResourceSegment(ctx context.Context, request *projpb.GetResourceSegmentRequest) (*projpb.ResourceSegment, error) {

	return t.data.GrpcClients.ProjProClient.XGetResourceSegment(ctx, &projpb.XGetResourceSegmentRequest{
		Id:     request.Id,
		UserId: "system",
	})
}
