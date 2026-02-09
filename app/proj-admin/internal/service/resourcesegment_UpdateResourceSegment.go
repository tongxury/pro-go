package service

import (
	"context"
	projpb "store/api/proj"
)

func (t ProjAdminService) UpdateResourceSegment(ctx context.Context, req *projpb.UpdateResourceSegmentRequest) (*projpb.ResourceSegment, error) {

	return t.data.GrpcClients.ProjProClient.XUpdateResourceSegment(ctx, &projpb.XUpdateResourceSegmentRequest{
		Id:     req.Id,
		Action: req.Action,
		RootId: req.RootId,
		UserId: "system",
	})
}
