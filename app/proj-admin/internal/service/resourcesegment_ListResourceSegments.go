package service

import (
	"context"
	projpb "store/api/proj"
)

func (t ProjAdminService) ListResourceSegments(ctx context.Context, params *projpb.ListResourceSegmentsRequest) (*projpb.ResourceSegmentList, error) {

	return t.data.GrpcClients.ProjProClient.XListResourceSegments(ctx, &projpb.XListResourceSegmentsRequest{
		Page:         params.Page,
		Keyword:      params.Keyword,
		SearchBy:     params.SearchBy,
		Category:     params.Category,
		Size:         params.Size,
		Status:       params.Status,
		Ids:          params.Ids,
		ReturnFields: params.ReturnFields,
		UserId:       "system",
	})
}
