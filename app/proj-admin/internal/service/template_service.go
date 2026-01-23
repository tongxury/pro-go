package service

import (
	"context"
	projpb "store/api/proj"
)

func (t ProjAdminService) UpdateTemplate_(ctx context.Context, request *projpb.UpdateTemplateRequest) (*projpb.Resource, error) {

	templates, err := t.data.GrpcClients.ProjProClient.XUpdateTemplate(ctx, request)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (t ProjAdminService) GetTemplate_(ctx context.Context, request *projpb.GetTemplateRequest) (*projpb.Resource, error) {

	templates, err := t.data.GrpcClients.ProjProClient.XGetTemplate(ctx, request)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (t ProjAdminService) ListTemplates_(ctx context.Context, request *projpb.ListTemplatesRequest) (*projpb.ResourceList, error) {

	templates, err := t.data.GrpcClients.ProjProClient.XListTemplates(ctx, &projpb.XListTemplatesRequest{
		Page:         request.Page,
		Category:     request.Category,
		Size:         request.Size,
		Ids:          request.Ids,
		UserId:       "system",
		ReturnFields: request.ReturnFields,
		Status:       request.Status,
	})
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (t ProjAdminService) AddTemplates_(ctx context.Context, request *projpb.AddTemplatesRequest) (*projpb.ResourceList, error) {

	templates, err := t.data.GrpcClients.ProjProClient.XAddTemplates(ctx, request)
	if err != nil {
		return nil, err
	}
	return templates, nil
}
