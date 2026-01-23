package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	databankpb "store/api/databank"
	"store/app/databank/internal/biz"
	"store/app/databank/internal/data"
	"store/pkg/krathelper"
)

type DatabankService struct {
	databankpb.UnimplementedDatabankServer
	userAssetBiz *biz.UserAssetBiz
	userFileBiz  *biz.UserFileBiz
	data         *data.Data
}

func NewDatabankService(userAssetBiz *biz.UserAssetBiz, userFileBiz *biz.UserFileBiz, data *data.Data) *DatabankService {
	return &DatabankService{
		userAssetBiz: userAssetBiz,
		userFileBiz:  userFileBiz,
		data:         data,
	}
}

func (t *DatabankService) ListAssets(ctx context.Context, params *databankpb.ListAssetsParams) (*databankpb.ListAssetsResult, error) {

	assets, err := t.userAssetBiz.ListAssets(ctx, params)
	if err != nil {
		log.Errorw("ListAssets err", err)
		return nil, err
	}
	return assets, nil
}

func (t *DatabankService) ListUserAssets(ctx context.Context, params *databankpb.ListUserAssetsParams) (*databankpb.ListAssetsResult, error) {

	userId := krathelper.RequireUserId(ctx)

	assets, err := t.userAssetBiz.ListUserAssets(ctx, userId, params)
	if err != nil {
		log.Errorw("ListUserAssets err", err)
		return nil, err
	}
	return assets, nil
}

func (t *DatabankService) AddAssets(ctx context.Context, params *databankpb.AddAssetsParams) (*emptypb.Empty, error) {

	userId := krathelper.RequireUserId(ctx)

	err := t.userAssetBiz.AddAssets(ctx, userId, params)
	if err != nil {
		log.Errorw("AddAssets err", err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (t *DatabankService) AddFiles(ctx context.Context, params *databankpb.AddFilesParams) (*databankpb.AddFilesResult, error) {

	userId := krathelper.RequireUserID(ctx)

	files, err := krathelper.GetMultipartFormFile(ctx, "file", "files")
	if err != nil {
		return nil, errors.BadRequest("", err.Error())
	}

	uploadedFiles, err := t.userFileBiz.AddFiles(ctx, userId, files, params.Public)
	if err != nil {
		log.Errorw("AddFiles err", err)
		return nil, err
	}

	return uploadedFiles, nil
}

//func (t *DatabankService) AddPubFiles(ctx context.Context, params *databankpb.AddFilesParams) (*databankpb.AddFilesResult, error) {
//
//	//formCodec := encoding.GetCodec("multipart/form-data")
//	files, err := kratosutil.GetMultipartFormFile(ctx, "files")
//	if err != nil {
//		return nil, errors.BadRequest("", err.Error())
//	}
//
//	uploadedFiles, err := t.userFileBiz.AddFiles(ctx, userId, files, params.Public)
//	if err != nil {
//		log.Errorw("AddFiles err", err)
//		return nil, err
//	}
//
//	return uploadedFiles, nil
//}
