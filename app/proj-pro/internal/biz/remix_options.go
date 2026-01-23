package biz

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/third/bytedance/volcengine"
	"time"
)

func (b *WorkflowBiz) initRemixCache() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// 立即执行一次
	b.UpdateRemixCache()

	for range ticker.C {
		b.UpdateRemixCache()
	}
}

func (b *WorkflowBiz) UpdateRemixCache() {
	ctx := context.Background()
	res := &projpb.RemixOptions{}

	// flowers
	flowers, err := b.data.Volcengine.GetFlowerList(ctx, volcengine.GetFlowerListParams{})
	if err == nil && flowers != nil {
		res.Flowers = helper.Mapping(*flowers, func(x volcengine.Flower) *projpb.RemixOptions_Flower {
			return &projpb.RemixOptions_Flower{
				MediaId: x.MediaId,
				Name:    x.Title,
				Cover:   x.Cover,
			}
		})
	}

	// fonts
	fonts, err := b.data.Volcengine.GetFontList(ctx, volcengine.GetFontListParams{})
	if err == nil && fonts != nil {
		res.Fonts = helper.Mapping(*fonts, func(x volcengine.Font) *projpb.RemixOptions_Font {
			return &projpb.RemixOptions_Font{
				MediaId: x.MediaId,
				Name:    x.Title,
			}
		})
	}

	// tones
	tones, err := b.data.Volcengine.GetToneList(ctx, volcengine.GetToneListParams{})
	if err == nil && tones != nil {
		res.Tones = helper.Mapping(*tones, func(x volcengine.Tone) *projpb.RemixOptions_Tone {
			return &projpb.RemixOptions_Tone{
				Id:          x.Id,
				Name:        x.Title,
				Description: x.Description,
				DownloadUrl: x.DownloadUrl,
			}
		})
	}

	_ = b.data.Cache.Remix.Set(ctx, res)
}

func (b *WorkflowBiz) GetRemixOptions(ctx context.Context) (*projpb.RemixOptions, error) {
	return b.data.Cache.Remix.Get(ctx)
}
