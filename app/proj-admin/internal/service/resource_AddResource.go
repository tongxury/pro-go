package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/videoz"
	"store/pkg/sdk/third/bytedance/tos"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

func (t ProjAdminService) AddItems(ctx context.Context, params *projpb.AddItemParams) (*projpb.ResourceList, error) {

	//newItem := &projpb.Item{
	//	Status:    "created",
	//	Title:     params.Title,
	//	Url:       params.Url,
	//	CoverUrl:  params.CoverUrl,
	//	CreatedAt: time.Now().Unix(),
	//}

	var items []*projpb.Resource
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
			}
		}

		items = append(items, &projpb.Resource{
			Status: "created",
			//Title:     x.Title,
			Url:       x.Url,
			CoverUrl:  x.CoverUrl,
			CreatedAt: time.Now().Unix(),
		})
	}

	//err := t.data.Elastics.CreateBulk(ctx, projpb.ESIndexItems, items)
	//if err != nil {
	//	return nil, err
	//}
	_, err := t.data.Mongo.Template.InsertMany(ctx, items...)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
