package service

import (
	"context"
	"fmt"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/third/bytedance/tos"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (t ProjAdminService) Mi6() {

	ctx := context.Background()

	list, err := t.data.Mongo.TemplateSegment.List(ctx, mgz.NotExists("extra.mi6"), options.Find().SetLimit(1))
	if err != nil {
		return
	}

	if len(list) == 0 {
		return
	}

	x := list[0]

	url, err := t.data.TOS.PutByUrl(ctx, tos.PutByUrlRequest{
		Bucket: "yoozyres",
		Url:    x.Root.Url,
		Suffix: ".mp4",
	})
	if err != nil {
		log.Errorw("err", err)
		return
	}

	_, err = t.data.Mongo.TemplateSegment.UpdateByIDIfExists(ctx, x.XId,
		mgz.Op().Set("root.url", url).Set("extra.mi6", true))
	if err != nil {
		log.Errorw("err", err)
		return
	}
}

func (t ProjAdminService) Mi4() {

	ctx := context.Background()

	list, err := t.data.Mongo.Template.List(ctx, mgz.NotExists("extra.mi4"), options.Find().SetLimit(1))
	if err != nil {
		return
	}

	if len(list) == 0 {
		return
	}

	x := list[0]

	url, err := t.data.TOS.PutByUrl(ctx, tos.PutByUrlRequest{
		Bucket: "yoozyres",
		Url:    x.CoverUrl,
		Suffix: ".jpg",
	})
	if err != nil {
		log.Errorw("err", err)
		return
	}

	_, err = t.data.Mongo.Template.UpdateByIDIfExists(ctx, x.XId,
		mgz.Op().Set("coverUrl", url))
	if err != nil {
		log.Errorw("err", err)
		return
	}

	_, err = t.data.Mongo.Template.UpdateByIDIfExists(ctx, x.XId,
		mgz.Op().Set("extra.mi4", true))
	if err != nil {
		log.Errorw("err", err)
		return
	}
}

func (t ProjAdminService) Mi2() {

	ctx := context.Background()

	list, err := t.data.Mongo.TemplateSegment.List(ctx, mgz.NotExists("extra.mi"), options.Find().SetLimit(1))
	if err != nil {
		return
	}

	if len(list) == 0 {
		return
	}

	x := list[0]

	url, err := t.data.TOS.PutByUrl(ctx, tos.PutByUrlRequest{
		Bucket: "yoozyres",
		Url:    x.Root.Url,
		Suffix: ".mp4",
	})
	if err != nil {
		log.Errorw("err", err)
		return
	}

	_, err = t.data.Mongo.TemplateSegment.UpdateByIDIfExists(ctx, x.XId,
		mgz.Op().Set("root.url", url))
	if err != nil {
		log.Errorw("err", err)
		return
	}

	for ii, xx := range x.HighlightFrames {
		url, err := t.data.TOS.PutByUrl(ctx, tos.PutByUrlRequest{
			Bucket: "yoozyres",
			Url:    xx.Url,
			Suffix: ".jpg",
		})
		if err != nil {
			log.Errorw("err", err)
			return
		}

		_, err = t.data.Mongo.TemplateSegment.UpdateByIDIfExists(ctx, x.XId,
			mgz.Op().Set(fmt.Sprintf("highlightFrames.%d.url", ii), url))
		if err != nil {
			log.Errorw("err", err)
			return
		}
	}

	_, err = t.data.Mongo.TemplateSegment.UpdateByIDIfExists(ctx, x.XId,
		mgz.Op().Set("extra.mi", true))
	if err != nil {
		log.Errorw("err", err)
		return
	}
}

func (t ProjAdminService) Mi3() {

	ctx := context.Background()

	list, err := t.data.Mongo.Template.List(ctx, mgz.NotExists("extra.mi"), options.Find().SetLimit(1))
	if err != nil {
		return
	}

	if len(list) == 0 {
		return
	}

	x := list[0]

	url, err := t.data.TOS.PutByUrl(ctx, tos.PutByUrlRequest{
		Bucket: "yoozyres",
		Url:    x.Url,
		Suffix: ".mp4",
	})
	if err != nil {
		log.Errorw("err", err)
		return
	}

	_, err = t.data.Mongo.Template.UpdateByIDIfExists(ctx, x.XId,
		mgz.Op().Set("url", url))
	if err != nil {
		log.Errorw("err", err)
		return
	}

	for i, xx := range x.Segments {

		for ii, xxx := range xx.HighlightFrames {
			url, err := t.data.TOS.PutByUrl(ctx, tos.PutByUrlRequest{
				Bucket: "yoozyres",
				Url:    xxx.Url,
				Suffix: ".jpg",
			})
			if err != nil {
				log.Errorw("err", err)
				return
			}

			_, err = t.data.Mongo.Template.UpdateByIDIfExists(ctx, x.XId,
				mgz.Op().Set(fmt.Sprintf("segments.%d.highlightFrames.%d.url", i, ii), url))
			if err != nil {
				log.Errorw("err", err)
				return
			}
		}

	}

	_, err = t.data.Mongo.Template.UpdateByIDIfExists(ctx, x.XId,
		mgz.Op().Set("extra.mi", true))
	if err != nil {
		log.Errorw("err", err)
		return
	}
}
