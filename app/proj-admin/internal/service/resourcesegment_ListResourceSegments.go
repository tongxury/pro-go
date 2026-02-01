package service

import (
	"context"
	"sort"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/third/bytedance/vikingdb"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjAdminService) ListResourceSegments(ctx context.Context, params *projpb.ListResourceSegmentsRequest) (*projpb.ResourceSegmentList, error) {

	size := helper.Select(params.Size > 0, params.Size, 24)
	page := helper.Select(params.Page > 0, params.Page, 1)

	var ids []string
	idSort := map[string]int{}
	if params.Keyword != "" {

		var err error
		var items *vikingdb.SearchResponse

		if params.SearchBy == "video" {
			items, err = t.data.VikingDB.SearchByKeywords(ctx, vikingdb.SearchByKeywordsRequest{
				SearchRequest: vikingdb.SearchRequest{
					CollectionName: "segment_video_coll",
					IndexName:      "segment_video_idx",
					Limit:          int(size),
				},
				Keywords: []string{params.Keyword},
			})
		} else {
			items, err = t.data.VikingDB.SearchByKeywords(ctx, vikingdb.SearchByKeywordsRequest{
				SearchRequest: vikingdb.SearchRequest{
					CollectionName: "segment_commodity_coll",
					IndexName:      "segment_commodity_idx",
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

		list, _, err := t.data.Mongo.TemplateSegment.ListAndCount(ctx,
			bson.M{"_id": bson.M{"$in": mgz.ObjectIds(ids)}},
			mgz.Find().
				Paging(0, 10).
				SetFields(params.ReturnFields).
				B(),
		)
		if err != nil {
			return nil, err
		}

		sort.Slice(list, func(i, j int) bool {
			return idSort[list[i].XId] < idSort[list[j].XId]
		})

		//for i := range list {
		//
		//	if len(list[i].HighlightFrames) > 0 {
		//		list[i].HighlightFrames = list[i].HighlightFrames[:1]
		//		list[i].HighlightFrames[0].Url = tos.Change(list[i].HighlightFrames[0].Url)
		//	}
		//}

		return &projpb.ResourceSegmentList{
			List: list,
			Page: 1,
			Size: size,
		}, nil
	}

	filter := bson.M{
		"status": bson.M{"$regex": "^completed"},
	}

	if params.Status != "" {
		filter["status"] = bson.M{"$regex": "^" + params.Status}
	}

	if len(params.Ids) > 0 {
		filter["_id"] = bson.M{"$in": mgz.ObjectIds(params.Ids)}
	}

	list, count, err := t.data.Mongo.TemplateSegment.ListAndCount(ctx,
		filter,
		mgz.Find().
			Paging(params.Page, params.Size).
			SetFields(params.ReturnFields).
			SetSort("createdAt", -1).
			B(),
	)
	if err != nil {
		return nil, err
	}

	for i := range list {
		list[i].Extra = &projpb.ResourceSegment_Extra{
			Polling: !strings.HasPrefix(list[i].Status, "completed"),
		}

		//if len(list[i].HighlightFrames) > 0 {
		//	list[i].HighlightFrames = list[i].HighlightFrames[:1]
		//	list[i].HighlightFrames[0].Url = tos.Change(list[i].HighlightFrames[0].Url)
		//}
	}

	return &projpb.ResourceSegmentList{
		List:  list,
		Total: count,
		Page:  page,
		Size:  size,
	}, nil

}
