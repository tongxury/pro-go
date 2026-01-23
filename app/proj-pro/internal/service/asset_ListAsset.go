package service

import (
	"context"
	"store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"store/pkg/sdk/third/bytedance/tos"
	strings "strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

func (t ProjService) ListAssets(ctx context.Context, params *projpb.ListAssetsRequest) (*projpb.AssetList, error) {

	userId := krathelper.RequireUserId(ctx)

	return t.XListAssets(ctx, &projpb.XListAssetsRequest{
		Page:         params.Page,
		Size:         params.Size,
		ReturnFields: params.ReturnFields,
		Status:       params.Status,
		UserId:       userId,
		Category:     params.Category,
	})

}
func (t ProjService) XListAssets(ctx context.Context, params *projpb.XListAssetsRequest) (*projpb.AssetList, error) {

	now := time.Now()

	filter := mgz.Filter()

	if params.Id != "" {
		filter = filter.EQ("_id", mgz.ObjectId(params.Id))
	}

	if params.UserId != "" {
		filter = filter.EQ("userId", params.UserId)
	}

	if params.Category != "" {
		filter = filter.EQ("category", params.Category)
	}

	if params.Status == "favorite" {
		filter = filter.EQ("favorite", true)
	} else {
		filter = filter.InString("status", strings.Split(params.Status, ","))
	}

	list, c, err := t.data.Mongo.Asset.ListAndCount(ctx,
		filter.B(),
		mgz.Find().
			Paging(params.Page, params.Size).
			SetFields("url,status,category,coverUrl,"+
				"commodity._id,commodity.coverUrl,commodity.images,"+
				"commodity.medias,createdAt,favorite,extra.context.status,extra.completedAt,"+
				"workflow._id,"+
				"remark").
			SetSort("createdAt", -1).
			B(),
	)
	if err != nil {
		return nil, err
	}

	for _, x := range list {
		//x.Extra = nil
		if x.CoverUrl == "" {
			if len(x.Commodity.GetMedias()) > 0 {
				x.CoverUrl = tos.Change(x.Commodity.GetMedias()[0].GetUrl())
			}
		}

		x.CoverUrl = tos.Change(x.CoverUrl)
		x.Url = tos.Change(x.Url)

		if len(x.Commodity.GetImages()) > 0 {
			x.Commodity.Images = tos.ChangeMany(x.Commodity.GetImages()[:1])
		}
	}

	//log.Debugw("params", params, "time", time.Since(now))

	// 补充workflow数据
	var workflowIds []string
	for _, x := range list {
		if x.Workflow.GetID() != "" {
			workflowIds = append(workflowIds, x.Workflow.GetID())
		}
	}

	if len(workflowIds) > 0 {
		workflows, err := t.data.Mongo.Workflow.List(ctx,
			mgz.Filter().Ids(workflowIds).B(),
			mgz.Find().SetFields("jobs,current,status,name").B(),
		)
		if err != nil {
			return nil, err
		}

		mp := mgz.ToMap(workflows)

		for _, x := range list {

			if x.Workflow.GetID() == "" {
				continue
			}

			workflow := mp[x.Workflow.GetID()]
			x.Workflow = workflow

			if x.CoverUrl == "" {
				images := workflow.GetDataBus().GetSegmentScript().GetImages()
				if len(images) > 0 {
					x.CoverUrl = tos.Change(images[0])
				}

				if len(workflow.GetDataBus().GetVideoGenerations()) > 0 {
					x.CoverUrl = tos.Change(workflow.GetDataBus().GetVideoGenerations()[0].CoverUrl)
				}
			}

			//if x.Title == "" {
			//	x.Title = workflow.GetDataBus().GetSegmentScript().GetScript()
			//}

			//x.Workflow.DataBus = nil
		}
	}

	log.Debugw("XListAssets", "", "cost", time.Since(now))

	return &projpb.AssetList{
		List:    list,
		Total:   c,
		Page:    params.Page,
		Size:    params.Size,
		HasMore: params.Size*params.Page < c,
	}, err
}
