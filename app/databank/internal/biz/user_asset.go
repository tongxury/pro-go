package biz

import (
	"context"
	"entgo.io/ent/dialect/sql"
	databankpb "store/api/databank"
	typepb "store/api/databank/types"
	enums "store/api/public/enums"
	"store/app/databank/internal/data"
	"store/app/databank/internal/data/repo/ent"
	"store/app/databank/internal/data/repo/ent/userasset"
	"store/pkg/sdk/conv"
)

type UserAssetBiz struct {
	data *data.Data
}

func NewUserAssetBiz(data *data.Data) *UserAssetBiz {
	return &UserAssetBiz{data: data}
}

func (t *UserAssetBiz) ListAssets(ctx context.Context, params *databankpb.ListAssetsParams) (*databankpb.ListAssetsResult, error) {

	q := t.data.Repos.EntClient.UserAsset.Query()

	if params.UserId != "" {
		q = q.Where(userasset.UserID(conv.Int64(params.UserId)))
	}

	count, err := q.Count(ctx)
	if err != nil {
		return nil, err
	}

	if params.Size > 0 && params.Page > 0 {
		q = q.Limit(int(params.Size)).Offset(int((params.Page - 1) * params.Size))
	}

	assets, err := q.All(ctx)
	if err != nil {
		return nil, err
	}

	return &databankpb.ListAssetsResult{
		Page:   params.Page,
		Size:   params.Size,
		Total:  int64(count),
		Assets: t.toPbAssets(ctx, assets),
	}, nil
}

func (t *UserAssetBiz) ListUserAssets(ctx context.Context, userId string, params *databankpb.ListUserAssetsParams) (*databankpb.ListAssetsResult, error) {

	q := t.data.Repos.EntClient.UserAsset.Query().
		Where(userasset.UserID(conv.Int64(userId)))

	count, err := q.Count(ctx)
	if err != nil {
		return nil, err
	}

	if params.Size > 0 && params.Page > 0 {
		q = q.Limit(int(params.Size)).Offset(int((params.Page - 1) * params.Size))
	}

	assets, err := q.All(ctx)
	if err != nil {
		return nil, err
	}

	return &databankpb.ListAssetsResult{
		Page:   params.Page,
		Size:   params.Size,
		Total:  int64(count),
		Assets: t.toPbAssets(ctx, assets),
	}, nil
}

func (t *UserAssetBiz) toPbAssets(ctx context.Context, assets []*ent.UserAsset) []*typepb.Asset {

	result := make([]*typepb.Asset, 0, len(assets))

	for _, x := range assets {
		y := typepb.Asset{
			Id: conv.String(x.ID),
			File: &typepb.File{
				Name: x.Name,
				Key:  x.Key,
			},
			School: x.School,
			Course: x.Course,
			Year:   x.Year,
			Status: &enums.Enum{
				Name:  x.Status.Name(),
				Value: x.Status.String(),
			},
			CreatedAt: x.CreatedAt.Unix(),
		}

		result = append(result, &y)
	}

	return result
}

func (t *UserAssetBiz) AddAssets(ctx context.Context, userId string, params *databankpb.AddAssetsParams) error {

	builders := make([]*ent.UserAssetCreate, 0, len(params.GetAssets()))
	for _, x := range params.GetAssets() {

		if x.File == nil {
			continue
		}

		builders = append(builders,
			t.data.Repos.EntClient.UserAsset.Create().
				SetUserID(conv.Int64(userId)).
				SetName(x.File.Name).
				SetKey(x.File.Key).
				SetSchool(x.School).
				SetCourse(x.Course).
				SetYear(x.Year),
		)
	}

	err := t.data.Repos.EntClient.UserAsset.CreateBulk(
		builders...,
	).OnConflict(sql.ConflictColumns(userasset.FieldKey)).
		DoNothing().Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}
