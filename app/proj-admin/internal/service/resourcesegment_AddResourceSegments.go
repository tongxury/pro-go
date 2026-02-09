package service

import (
	"context"
	projpb "store/api/proj"
)

func (t ProjAdminService) AddResourceSegments(ctx context.Context, params *projpb.AddResourcesSegmentRequest) (*projpb.ResourceSegmentList, error) {

	var items []*projpb.XAddResourcesSegmentRequest_Item
	for _, x := range params.GetItems() {
		items = append(items, &projpb.XAddResourcesSegmentRequest_Item{
			Url:       x.Url,
			CoverUrl:  x.CoverUrl,
			TimeStart: x.TimeStart,
			TimeEnd:   x.TimeEnd,
		})
	}

	return t.data.GrpcClients.ProjProClient.XAddResourceSegments(ctx, &projpb.XAddResourcesSegmentRequest{
		Items:  items,
		UserId: "system",
	})
}
