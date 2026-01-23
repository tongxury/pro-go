package service

import (
	"context"
	projpb "store/api/proj"
)

func (t ProjService) CreateTemplate(ctx context.Context, request *projpb.CreateTemplateRequest) (*projpb.Resource, error) {
	return &projpb.Resource{}, nil
}
