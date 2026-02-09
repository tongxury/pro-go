package service

import (
	"context"
	"sort"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/third/bytedance/vikingdb"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjService) ListResourceSegments(ctx context.Context, params *projpb.ListResourceSegmentsRequest) (*projpb.ResourceSegmentList, error) {

	// For public API, we might use userId from context for "collected" status check
	userId, _ := krathelper.GetUserId(ctx)

	size := helper.Select(params.Size > 0, params.Size, 24)
	page := helper.Select(params.Page > 0, params.Page, 1)

	var finalIds []string
	useIds := false

	// 1. Gather Candidate IDs from Collected
	if params.Collected {
		useIds = true
		if userId == "" {
			return nil, errors.Unauthorized("unauthorized", "login required")
		}
		list, err := t.data.Mongo.ResourceSegmentCollection.List(ctx,
			bson.M{"user._id": userId},
			mgz.Find().SetSort("createdAt", -1).Limit(2000).B(),
		)
		if err != nil {
			return nil, err
		}
		if len(list) == 0 {
			return &projpb.ResourceSegmentList{Page: page, Size: size}, nil
		}
		for _, item := range list {
			seg := item.ResourceSegment
			if seg == nil {
				seg = item.ResourceSemgentLegacy
			}
			if seg != nil {
				finalIds = append(finalIds, seg.XId)
			}
		}
	}

	// 2. Gather Candidate IDs from Keyword Search (and Intersect)
	if params.Keyword != "" {
		var err error
		var items *vikingdb.SearchResponse
		limit := helper.Select(useIds, 500, int(size))

		searchReq := vikingdb.SearchByKeywordsRequest{
			SearchRequest: vikingdb.SearchRequest{
				CollectionName: helper.Select(params.SearchBy == "video", "segment_video_coll", "segment_commodity_coll"),
				IndexName:      helper.Select(params.SearchBy == "video", "segment_video_idx", "segment_commodity_idx"),
				Limit:          limit,
			},
			Keywords: []string{params.Keyword},
		}

		items, err = t.data.VikingDB.SearchByKeywords(ctx, searchReq)
		if err != nil {
			return nil, err
		}

		var searchIds []string
		for _, item := range items.Data {
			searchIds = append(searchIds, item.Id)
		}

		if useIds {
			// Intersection
			collectedMap := make(map[string]bool)
			for _, id := range finalIds {
				collectedMap[id] = true
			}
			var intersected []string
			for _, id := range searchIds {
				if collectedMap[id] {
					intersected = append(intersected, id)
				}
			}
			finalIds = intersected
		} else {
			finalIds = searchIds
			useIds = true
		}

		if len(finalIds) == 0 {
			return &projpb.ResourceSegmentList{Page: page, Size: size}, nil
		}
	}

	// 3. Final Fetching and Post-processing
	var segments []*projpb.ResourceSegment
	var total int64
	var err error

	if useIds {
		total = int64(len(finalIds))
		start := (page - 1) * size
		if start >= total {
			return &projpb.ResourceSegmentList{Total: total, Page: page, Size: size}, nil
		}
		end := helper.Select(start+size > total, total, start+size)
		pagedIds := finalIds[start:end]

		idOrder := map[string]int{}
		for i, id := range pagedIds {
			idOrder[id] = i
		}

		segments, _, err = t.data.Mongo.TemplateSegment.ListAndCount(ctx,
			bson.M{"_id": bson.M{"$in": mgz.ObjectIds(pagedIds)}},
			mgz.Find().SetFields(params.ReturnFields).B(),
		)
		if err != nil {
			return nil, err
		}
		sort.Slice(segments, func(i, j int) bool {
			return idOrder[segments[i].XId] < idOrder[segments[j].XId]
		})
	} else {
		filter := bson.M{"status": bson.M{"$regex": "^completed"}}
		if params.Status != "" {
			filter["status"] = bson.M{"$regex": "^" + params.Status}
		}
		if len(params.Ids) > 0 {
			filter["_id"] = bson.M{"$in": mgz.ObjectIds(params.Ids)}
		}

		segments, total, err = t.data.Mongo.TemplateSegment.ListAndCount(ctx,
			filter,
			mgz.Find().Paging(params.Page, params.Size).SetFields(params.ReturnFields).SetSort("createdAt", -1).B(),
		)
		if err != nil {
			return nil, err
		}
	}

	// 4. Common Post-processing
	for i := range segments {
		segments[i].Extra = &projpb.ResourceSegment_Extra{
			Polling: !strings.HasPrefix(segments[i].Status, "completed"),
		}
	}

	if userId != "" {
		t.fillCollectionStatus(ctx, userId, segments)
	}

	return &projpb.ResourceSegmentList{
		List:  segments,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

func (t ProjService) fillCollectionStatus(ctx context.Context, userId string, list []*projpb.ResourceSegment) {
	if len(list) == 0 {
		return
	}
	var segmentIds []string
	for _, item := range list {
		segmentIds = append(segmentIds, item.XId)
	}

	collections, err := t.data.Mongo.ResourceSegmentCollection.List(ctx, bson.M{
		"user._id": userId,
		"$or": []bson.M{
			{"resourceSegment._id": bson.M{"$in": segmentIds}},
			{"resourceSemgent._id": bson.M{"$in": segmentIds}},
		},
	})
	if err != nil {
		return
	}

	collectedMap := make(map[string]bool)
	for _, c := range collections {
		if c.ResourceSegment != nil {
			collectedMap[c.ResourceSegment.XId] = true
		} else if c.ResourceSemgentLegacy != nil {
			collectedMap[c.ResourceSemgentLegacy.XId] = true
		}
	}

	for _, item := range list {
		if collectedMap[item.XId] {
			item.Collected = true
		}
	}
}
