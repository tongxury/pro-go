package service

import (
	"context"
	projpb "store/api/proj"
	"strings"
)

func (t *ProjService) ListRemixOptions(ctx context.Context, req *projpb.ListRemixOptionsRequest) (*projpb.RemixOptions, error) {

	cache, err := t.workflow.GetRemixOptions(ctx)
	if err != nil {
		return nil, err
	}

	if cache == nil {
		return &projpb.RemixOptions{}, nil
	}

	fields := req.Fields
	if fields == "" {
		fields = "flowers,fonts,tones"
	}

	res := &projpb.RemixOptions{}

	if strings.Contains(fields, "flowers") {
		res.Flowers = cache.Flowers
	}

	if strings.Contains(fields, "fonts") {
		res.Fonts = cache.Fonts
	}

	if strings.Contains(fields, "tones") {
		res.Tones = cache.Tones
	}

	return res, nil
}
