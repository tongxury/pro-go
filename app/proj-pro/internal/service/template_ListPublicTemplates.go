package service

import (
	"context"
	projpb "store/api/proj"
)

func (t ProjService) ListPublicTemplates(ctx context.Context, params *projpb.ListPublicTemplatesRequest) (*projpb.ResourceList, error) {

	return t.XListTemplates(ctx, &projpb.XListTemplatesRequest{
		Page:         params.Page,
		Category:     params.Category,
		Size:         params.Size,
		Ids:          params.Ids,
		UserId:       "system",
		ReturnFields: params.ReturnFields,
		Status:       params.Status,
	})
}
