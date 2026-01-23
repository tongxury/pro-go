package service

import (
	"context"
	"sort"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/third/bytedance/tos"
	"store/pkg/sdk/third/bytedance/vikingdb"
	"strings"
)

func (t ProjService) ListTemplates(ctx context.Context, params *projpb.ListTemplatesRequest) (*projpb.ResourceList, error) {

	userId := krathelper.RequireUserId(ctx)

	return t.XListTemplates(ctx, &projpb.XListTemplatesRequest{
		Page:         params.Page,
		Category:     params.Category,
		Size:         params.Size,
		Ids:          params.Ids,
		UserId:       userId,
		ReturnFields: params.ReturnFields,
		Status:       params.Status,
	})
}

func (t ProjService) extra(list []*projpb.Resource) {

	for i := range list {

		//var polling bool = true
		//if len(list[i].Segments) > 0 {
		//	//polling = len(list[i].Segments[0].HighlightFrames) == 0
		//}

		list[i].CoverUrl = tos.Change(list[i].CoverUrl)
		list[i].Url = tos.Change(list[i].Url)

		list[i].Extra = &projpb.Resource_Extra{
			Polling: list[i].Status != "completed",
		}
	}

}

func (t ProjService) XListTemplates(ctx context.Context, params *projpb.XListTemplatesRequest) (*projpb.ResourceList, error) {

	size := helper.Select(params.Size > 0, params.Size, 20)
	page := helper.Select(params.Page > 0, params.Page, 1)

	if params.Keyword != "" {
		return t.searchTemplates(ctx, params)
	}

	filter := mgz.Filter()

	if params.UserId != "" {
		filter = filter.EQ("userId", params.UserId)
	}

	if len(params.Ids) > 0 {
		filter = filter.Ids(params.Ids)
	}

	if params.Status != "" {
		filter = filter.EQ("status", params.Status)
	}

	if params.ReturnFields != "" && !strings.Contains(params.ReturnFields, "status") {
		params.ReturnFields += ",status"
	}

	list, count, err := t.data.Mongo.Template.ListAndCount(ctx,
		filter.B(),
		mgz.Find().
			SetFields(params.ReturnFields).
			PageSize(page, size).
			SetSort("createdAt", -1).
			B(),
	)
	if err != nil {
		return nil, err
	}

	t.extra(list)

	return &projpb.ResourceList{
		List:  list,
		Total: count,
		Page:  page,
		Size:  size,
	}, nil

}

func (t ProjService) searchTemplates(ctx context.Context, params *projpb.XListTemplatesRequest) (*projpb.ResourceList, error) {

	size := helper.Select(params.Size > 0, params.Size, 20)
	if params.UserId == "" {
		params.UserId = "system"
	}

	var ids []string
	idSort := map[string]int{}

	if params.Keyword != "" {

		var err error
		var items *vikingdb.SearchResponse

		items, err = t.data.VikingDB.SearchByMultiModal(ctx, vikingdb.SearchByMultiModalRequest{
			SearchRequest: vikingdb.SearchRequest{
				CollectionName: "template_commodity_coll",
				IndexName:      "template_commodity_text_userId_idx",
				Limit:          int(size),
				Filter: map[string]interface{}{
					"op":    "must",
					"field": "userId",
					"conds": []string{params.UserId},
				},
			},
			Text:            params.Keyword,
			NeedInstruction: false,
		})

		if err != nil {
			return nil, err
		}

		if len(items.Data) == 0 {
			return nil, nil
		}

		for i, item := range items.Data {

			if helper.InSlice(item.Id, ids) {
				continue
			}

			ids = append(ids, item.Id)
			idSort[item.Id] = i
		}

	}

	filters := mgz.Filter()

	if len(ids) > 0 {
		filters.Ids(ids)
	}

	if params.UserId != "" {
		filters.EQ("userId", params.UserId)
	}

	list, _, err := t.data.Mongo.Template.ListAndCount(ctx,
		filters.B(),
		mgz.Find().Limit(size).B(),
	)
	if err != nil {
		return nil, err
	}

	if len(ids) > 0 {
		sort.Slice(list, func(i, j int) bool {
			return idSort[list[i].XId] < idSort[list[j].XId]
		})
	}

	t.extra(list)

	return &projpb.ResourceList{
		List: list,
		Page: 1,
		Size: size,
	}, nil
}
