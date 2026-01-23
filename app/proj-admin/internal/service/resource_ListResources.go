package service

import (
	"context"
	"sort"
	projpb "store/api/proj"
	"store/app/proj-admin/internal/biz"
	"store/app/proj-admin/internal/data"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/third/bytedance/vikingdb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProjAdminService struct {
	projpb.UnimplementedProjAdminServiceServer
	data *data.Data
	item *biz.ItemBiz
}

func NewProjAdminService(data *data.Data, item *biz.ItemBiz) *ProjAdminService {
	return &ProjAdminService{
		data: data,
		item: item,
	}
}

func (t ProjAdminService) ListItems(ctx context.Context, params *projpb.ListItemsParams) (*projpb.ResourceList, error) {

	size := helper.Select(params.Size > 0, params.Size, 20)
	page := helper.Select(params.Page > 0, params.Page, 1)

	var ids []string
	idSort := map[string]int{}

	if params.Keyword != "" {

		var err error
		var items *vikingdb.SearchResponse

		if params.SearchBy == "video" {
			items, err = t.data.VikingDB.SearchByKeywords(ctx, vikingdb.SearchByKeywordsRequest{
				SearchRequest: vikingdb.SearchRequest{
					CollectionName: "template_commodity_coll",
					IndexName:      "template_commodity_text_userId_idx",
					Limit:          int(size),
				},
				Keywords: []string{params.Keyword},
			})
		} else {

			items, err = t.data.VikingDB.SearchByKeywords(ctx, vikingdb.SearchByKeywordsRequest{
				SearchRequest: vikingdb.SearchRequest{
					CollectionName: "template_commodity_coll",
					IndexName:      "template_commodity_text_userId_idx",
					Limit:          int(size),
				},
				Keywords: []string{params.Keyword},
			})
		}

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

		list, _, err := t.data.Mongo.Template.ListAndCount(ctx,
			bson.M{"_id": bson.M{"$in": mgz.ObjectIds(ids)}},
			//options.Find().SetLimit(size).SetSkip(size*(page-1)),
		)
		if err != nil {
			return nil, err
		}

		sort.Slice(list, func(i, j int) bool {
			return idSort[list[i].XId] < idSort[list[j].XId]
		})

		return &projpb.ResourceList{
			List: list,
			Page: 1,
			Size: size,
		}, nil
	}

	filter := bson.M{}

	if len(ids) > 0 {
		filter["_id"] = bson.M{"$in": mgz.ObjectIds(ids)}
	}

	list, count, err := t.data.Mongo.Template.ListAndCount(ctx,
		filter,
		options.Find().SetLimit(size).SetSkip(size*(page-1)).
			SetSort(bson.M{"createdAt": -1}),
	)
	if err != nil {
		return nil, err
	}

	for i := range list {

		var polling bool = true
		if len(list[i].Segments) > 0 {
			polling = len(list[i].Segments[0].HighlightFrames) == 0
		}

		list[i].Extra = &projpb.Resource_Extra{
			Polling: polling,
		}
	}

	return &projpb.ResourceList{
		List:  list,
		Total: count,
		Page:  page,
		Size:  size,
	}, nil

}

func (t ProjAdminService) PutItem(ctx context.Context, params *projpb.Item) (*projpb.Item, error) {

	//err := t.data.Elastics.Replace(ctx, projpb.ESIndexItems, params.Id, params)
	//if err != nil {
	//	return nil, err
	//}

	return nil, nil
}

func (t ProjAdminService) UpdateItemTags(ctx context.Context, params *projpb.UpdateItemTagParams) (*projpb.Item, error) {

	//err := t.data.Elastics.UpdateFields(ctx, projpb.ESIndexItems, params.Id, map[string]interface{}{
	//	"tags": params.Tags,
	//})
	//if err != nil {
	//	return nil, err
	//}

	return nil, nil
}

func (t ProjAdminService) UpdateItemCategory(ctx context.Context, params *projpb.UpdateItemCategoryParams) (*projpb.Item, error) {

	//err := t.data.Elastics.UpdateFields(ctx, projpb.ESIndexItems, params.Id, map[string]interface{}{
	//	"commodity": map[string]interface{}{
	//		"category": params.Category,
	//	},
	//})
	//if err != nil {
	//	return nil, err
	//}

	return nil, nil
}

func (t ProjAdminService) GetItem(ctx context.Context, params *projpb.GetItemParams) (*projpb.Resource, error) {

	id, err := t.data.Mongo.Template.FindByID(ctx, params.Id)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (t ProjAdminService) UpdateItem(ctx context.Context, params *projpb.UpdateItemParams) (*projpb.Resource, error) {

	var err error
	switch params.Action {
	case "delete":
		//err = t.data.Elastics.Delete(ctx, projpb.ESIndexItems, params.Id)
		err = t.data.Mongo.Template.DeleteByID(ctx, params.Id)
	case "refresh":

		_, err := t.data.Mongo.TemplateSegment.Delete(ctx, bson.M{"root.id": params.Id})
		if err != nil {
			return nil, err
		}

		_, err = t.data.Mongo.Template.UpdateByIDXX(ctx, params.Id, mgz.Set(bson.M{"status": "created", "segments": nil}))

		//f, err := t.data.Elastics.DeleteByRequest(ctx, projpb.ESIndexItemSegments, deletebyquery.Request{
		//	Query: elastics.NewTermQuery("root.id.keyword", params.Id),
		//})
		//
		//fmt.Println(f, err)
		//
		//err = t.data.Elastics.UpdateFields(ctx, projpb.ESIndexItems, params.Id, map[string]interface{}{
		//	"status": "created",
		//})

	}

	if err != nil {
		return nil, err
	}

	return nil, nil
}
