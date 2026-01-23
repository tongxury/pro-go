package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/sdk/third/bytedance/vikingdb"
)

func (t ProjAdminService) SearchMatchedItems(ctx context.Context, params *projpb.SearchMatchedItemsParams) (*projpb.ResourceList, error) {

	items, err := t.data.VikingDB.SearchByKeywords(ctx, vikingdb.SearchByKeywordsRequest{
		SearchRequest: vikingdb.SearchRequest{
			CollectionName: "template_commodity_coll",
			IndexName:      "template_commodity_text_userId_idx",
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

	tpl, err := t.data.Mongo.Template.FindByID(ctx, items.Data[0].Id)
	if err != nil {
		return nil, err
	}

	return &projpb.ResourceList{
		List: []*projpb.Resource{tpl},
	}, nil
}
