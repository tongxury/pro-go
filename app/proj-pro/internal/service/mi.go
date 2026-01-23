package service

import (
	"context"
	"fmt"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/videoz"
	"store/pkg/sdk/third/bytedance/tos"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (t ProjService) AssetCover() {

	ctx := context.Background()

	list, err := t.data.Mongo.Asset.List(ctx,
		mgz.Filter().
			EQ("status", "completed").
			NotExists("extra.mcover1").
			B(),
	)
	if err != nil {
		return
	}

	for _, x := range list {

		coverFrame, err := videoz.GetFrameByUrl(x.Url, 1)
		if err != nil {
			log.Errorw(" err", err, "url", x.XId)
			return
		}
		coverUrl, err := t.data.TOS.Put(ctx, tos.PutRequest{
			Bucket:  "yoozyres",
			Content: coverFrame,
			Key:     helper.MD5(coverFrame) + ".jpg",
		})
		if err != nil {
			log.Errorw(" err", err, "url", x.XId)
			return
		}

		_, err = t.data.Mongo.Asset.UpdateByIDXX(ctx, x.XId, mgz.Set(bson.M{
			"coverUrl":      coverUrl,
			"extra.mcover1": true,
		}))
		if err != nil {
			return
		}

	}
}

func (t ProjService) Mi7() {

	ctx := context.Background()

	list, err := t.data.Mongo.Asset.List(ctx, mgz.NotExists("extra.mi7"))
	if err != nil {
		return
	}

	for _, x := range list {

		segment, err := t.data.GrpcClients.ProjAdminClient.GetResourceSegment(ctx, &projpb.GetResourceSegmentRequest{
			Id: x.Segment.XId,
		})
		if err != nil {
			log.Errorw("GetResourceSegment err", err, "x_id", x.Segment.XId)
			continue
		}

		if segment == nil {
			_, err = t.data.Mongo.Asset.UpdateByIDIfExists(ctx, x.XId,
				mgz.Op().Set("extra.mi7", true))
			if err != nil {
				return
			}

			return
		}

		_, err = t.data.Mongo.Asset.UpdateByIDIfExists(ctx, x.XId,
			mgz.Op().Set("segment", segment).Set("extra.mi7", true))
		if err != nil {
			return
		}
	}

}

func (t ProjService) Mi3() {

	ctx := context.Background()

	list, err := t.data.Mongo.Commodity.List(ctx, mgz.NotExists("extra.mi"), options.Find().SetLimit(1))
	if err != nil {
		return
	}

	if len(list) == 0 {
		return
	}

	x := list[0]

	for i, xx := range x.Medias {

		url, err := t.data.TOS.PutByUrl(ctx, tos.PutByUrlRequest{
			Bucket: "yoozyres",
			Url:    xx.Url,
			Suffix: ".jpg",
		})
		if err != nil {
			return
		}

		_, err = t.data.Mongo.Commodity.UpdateByIDIfExists(ctx, x.XId,
			mgz.Op().Set(fmt.Sprintf("medias.%d.url", i), url))
		if err != nil {
			return
		}
	}

	_, err = t.data.Mongo.Commodity.UpdateByIDIfExists(ctx, x.XId,
		mgz.Op().Set("extra.mi", true))
	if err != nil {
		return
	}
}

func (t ProjService) Mi2() {

	ctx := context.Background()

	list, err := t.data.Mongo.Commodity.List(ctx, mgz.NotExists("extra.mi"), options.Find().SetLimit(1))
	if err != nil {
		return
	}

	if len(list) == 0 {
		return
	}

	x := list[0]

	for i, xx := range x.Medias {

		url, err := t.data.TOS.PutByUrl(ctx, tos.PutByUrlRequest{
			Bucket: "yoozyres",
			Url:    xx.Url,
			Suffix: ".jpg",
		})
		if err != nil {
			return
		}

		_, err = t.data.Mongo.Commodity.UpdateByIDIfExists(ctx, x.XId,
			mgz.Op().Set(fmt.Sprintf("medias.%d.url", i), url))
		if err != nil {
			return
		}
	}

	_, err = t.data.Mongo.Commodity.UpdateByIDIfExists(ctx, x.XId,
		mgz.Op().Set("extra.mi", true))
	if err != nil {
		return
	}
}

func (t ProjService) Mi1() {

	ctx := context.Background()

	list, err := t.data.Mongo.TaskSegment.List(ctx, mgz.NotExists("extra.mi2"), options.Find().SetLimit(1))
	if err != nil {
		return
	}

	if len(list) == 0 {
		return
	}

	x := list[0]

	url, err := t.data.TOS.PutByUrl(ctx, tos.PutByUrlRequest{
		Bucket: "yoozyres",
		//Url:    x.,
		Suffix: ".mp4",
	})
	if err != nil {
		return
	}

	_, err = t.data.Mongo.Asset.UpdateByIDIfExists(ctx, x.XId,
		mgz.Op().Set("url", url).Set("extra.mi", true))
	if err != nil {
		return
	}

}

func (t ProjService) Mi() {

	ctx := context.Background()

	list, err := t.data.Mongo.Asset.List(ctx, mgz.NotExists("extra.mi2"), options.Find().SetLimit(1))
	if err != nil {
		return
	}

	if len(list) == 0 {
		return
	}

	x := list[0]

	url, err := t.data.TOS.PutByUrl(ctx, tos.PutByUrlRequest{
		Bucket: "yoozyres",
		//Url:    x.,
		Suffix: ".mp4",
	})
	if err != nil {
		return
	}

	_, err = t.data.Mongo.Asset.UpdateByIDIfExists(ctx, x.XId,
		mgz.Op().Set("url", url).Set("extra.mi", true))
	if err != nil {
		return
	}

}

func (t ProjService) Mi5() {

	ctx := context.Background()

	list, err := t.data.Mongo.Asset.List(ctx, mgz.NotExists("extra.mi6"), options.Find().SetLimit(1))
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
		return
	}

	_, err = t.data.Mongo.Asset.UpdateByIDIfExists(ctx, x.XId,
		mgz.Op().Set("coverUrl", url).Set("extra.mi6", true))
	if err != nil {
		return
	}

}
