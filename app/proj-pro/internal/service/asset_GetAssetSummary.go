package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"store/pkg/sdk/helper/wg"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (t ProjService) GetAssetSummary(ctx context.Context, request *emptypb.Empty) (*projpb.AssetSummary, error) {

	userId := krathelper.RequireUserId(ctx)

	var (
		total      int64
		favorite   int64
		processing int64
		completed  int64
		failed     int64
	)

	wg.WaitGroupFunctions(ctx,
		func(ctx context.Context) error {
			total, _ = t.data.Mongo.Asset.Count(ctx, mgz.Filter().EQ("userId", userId).B())
			return nil
		},
		func(ctx context.Context) error {
			favorite, _ = t.data.Mongo.Asset.Count(ctx, mgz.Filter().EQ("userId", userId).EQ("favorite", true).B())
			return nil
		},
		func(ctx context.Context) error {
			processing, _ = t.data.Mongo.Asset.Count(ctx, mgz.Filter().EQ("userId", userId).NotEQ("status", "completed").B())
			return nil
		},
		func(ctx context.Context) error {
			completed, _ = t.data.Mongo.Asset.Count(ctx, mgz.Filter().EQ("userId", userId).EQ("status", "completed").B())
			return nil
		},
		func(ctx context.Context) error {
			failed, _ = t.data.Mongo.Asset.Count(ctx, mgz.Filter().EQ("userId", userId).EQ("status", "failed").B())
			return nil
		},
	)

	return &projpb.AssetSummary{
		Total:      total,
		Favorite:   favorite,
		Processing: processing,
		Completed:  completed,
		Failed:     failed,
	}, nil
}
