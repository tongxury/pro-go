package service

import (
	"context"
	projpb "store/api/proj"
	usercenter "store/api/usercenter"
	"store/pkg/krathelper"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/videoz"
	"store/pkg/sdk/third/bytedance/tos"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

func (t ProjService) AddTemplates(ctx context.Context, params *projpb.AddTemplatesRequest) (*projpb.ResourceList, error) {

	//newItem := &projpb.Item{
	//	Status:    "created",
	//	Title:     params.Title,
	//	Url:       params.Url,
	//	CoverUrl:  params.CoverUrl,
	//	CreatedAt: time.Now().Unix(),
	//}

	userId := krathelper.RequireUserId(ctx)

	if len(params.Items) == 0 {
		return nil, errors.BadRequest("emptyItems", "")
	}

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
			User:      &usercenter.User{XId: userId},
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

func (t ProjService) XAddTemplates(ctx context.Context, params *projpb.AddTemplatesRequest) (*projpb.ResourceList, error) {

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
			User:      &usercenter.User{XId: "system"},
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

//rpc _ListTemplates(ListTemplatesRequest) returns (api.proj.ResourceList) {}
//rpc _AddTemplates(AddTemplatesRequest) returns (api.proj.ResourceList) {}
//rpc _GetTemplate(GetTemplateRequest) returns (api.proj.Resource) {}
//rpc _UpdateTemplate(UpdateTemplateRequest) returns (api.proj.Resource) {}
