package service

import (
	"context"
	projpb "store/api/proj"
	ucpb "store/api/usercenter"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/videoz"
	"store/pkg/sdk/third/bytedance/tos"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

func (t ProjService) XAddResourceSegments(ctx context.Context, params *projpb.XAddResourcesSegmentRequest) (*projpb.ResourceSegmentList, error) {

	var items []*projpb.ResourceSegment
	for _, x := range params.GetItems() {

		if x.CoverUrl == "" {
			bytes, err := videoz.GetFrameByUrl(x.Url, 1)
			if err == nil {

				x.CoverUrl, err = t.data.TOS.Put(ctx, tos.PutRequest{
					Bucket:  "yoozyres",
					Content: bytes,
					Key:     helper.MD5(bytes) + ".jpg",
				})
				if err != nil {
					return nil, err
				}
			} else {
				log.Errorw("GetFrameByUrl err", err)
				return nil, err
			}
		}

		status := "processing_created"
		if x.TimeEnd > 0 {
			status = "processing_segmented"
		}

		y := &projpb.ResourceSegment{
			Status: status,
			//Title:     x.Title,
			Root: &projpb.Resource{
				Url:      x.Url,
				CoverUrl: x.CoverUrl,
				User:     &ucpb.User{XId: params.UserId},
			},

			TimeStart: x.TimeStart,
			TimeEnd:   x.TimeEnd,

			CreatedAt: time.Now().Unix(),
		}

		items = append(items, y)
	}

	_, err := t.data.Mongo.TemplateSegment.InsertMany(ctx, items...)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
