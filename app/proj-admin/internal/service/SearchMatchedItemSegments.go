package service

import (
	"context"
	"encoding/json"
	projpb "store/api/proj"
	"store/pkg/sdk/third/bytedance/vikingdb"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

func (t ProjAdminService) SearchMatchedItemSegments(ctx context.Context, params *projpb.SearchMatchedItemSegmentsParams) (*projpb.ResourceSegmentList, error) {

	items, err := t.data.VikingDB.SearchByKeywords(ctx, vikingdb.SearchByKeywordsRequest{
		SearchRequest: vikingdb.SearchRequest{
			CollectionName: "segment_commodity_coll",
			IndexName:      "segment_commodity_idx",
			Limit:          int(params.Limit),
		},
		Keywords: []string{params.Keyword},
	})
	if err != nil {
		return nil, err
	}

	if len(items.Data) == 0 {
		return nil, nil
	}

	tpl, err := t.data.Mongo.TemplateSegment.FindByID(ctx, items.Data[0].Id)
	if err != nil {
		return nil, err
	}

	return &projpb.ResourceSegmentList{
		List: []*projpb.ResourceSegment{tpl},
	}, nil
}

func (t ProjAdminService) parseToItemSegments(hits []types.Hit) []*projpb.ItemSegment {
	var items []*projpb.ItemSegment
	for _, x := range hits {

		var y *projpb.ItemSegment
		_ = json.Unmarshal(x.Source_, &y)

		//y.Id = *x.Id_

		if y.Commodity != nil {

			//y.Commodity.CategoryVector = nil
		}

		items = append(items, y)
	}

	return items
}
