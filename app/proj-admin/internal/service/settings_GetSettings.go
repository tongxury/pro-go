package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjAdminService) GetSettings(ctx context.Context, req *projpb.GetSettingsRequest) (*projpb.AppSettings, error) {

	// 返回
	if req.Fields == "prompts" {

		//t.data.Mongo.Settings.UpdateByIDIfExists(ctx, "6928a405029aef14ece37500",
		//	mgz.Op().Set(
		//		"prompts",
		//		map[string]*projpb.Prompt{
		//			"demoq": {
		//				Content:   "demo content",
		//				UpdatedAt: time.Now().Unix(),
		//			},
		//		},
		//	),
		//)

		rs, err := t.data.Mongo.Settings.Find(ctx, bson.M{})
		if err != nil {
			return nil, err
		}
		if len(rs) == 0 {
			return &projpb.AppSettings{}, nil
		}

		for _, x := range rs[0].Prompts {
			x.Content = ""
		}

		return rs[0], nil

	}

	rs, err := t.data.Mongo.Settings.Find(ctx, bson.M{}, mgz.Find().SetFields(req.Fields).B())
	if err != nil {
		return nil, err
	}

	//if len(rs) == 0 {
	//	r, err := t.data.Mongo.Settings.Insert(ctx, &projpb.AppSettings{
	//		Prompts: map[string]*projpb.Prompt{
	//			"demo": {
	//				Content:   "demo content",
	//				UpdatedAt: time.Now().Unix(),
	//			},
	//		},
	//		VideoTemplate: &projpb.Prompt{
	//			Content: "xx",
	//		},
	//		VideoHighlight: &projpb.Prompt{
	//			Content: "xx",
	//		},
	//	})
	//
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	return r, nil
	//}

	return rs[0], nil
}
