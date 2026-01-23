package biz

import (
	"context"
	"fmt"
	databankpb "store/api/databank"
	typepb "store/api/databank/types"
	"store/app/databank/internal/data"
	"store/app/databank/internal/data/repo/ent"
	"store/pkg/krathelper"
	"store/pkg/sdk/helper/filed"
)

type UserFileBiz struct {
	data *data.Data
}

func NewUserFileBiz(data *data.Data) *UserFileBiz {
	return &UserFileBiz{data: data}
}

func (t *UserFileBiz) AddFiles(ctx context.Context, userID int64, files krathelper.FormFiles, public bool) (*databankpb.AddFilesResult, error) {

	if len(files) == 0 {
		return nil, nil
	}

	// 上传
	//var dataFiles data.Files
	//for _, x := range files {
	//	dataFiles = append(dataFiles, data.File{
	//		Name: x.Filename,
	//		Body: x.Body,
	//		Md5:  helper.MD5(x.Body),
	//	})
	//}

	uploadedFiles, err := t.data.Uploader.Upload(ctx, files, fmt.Sprintf("user-files/%d", userID))
	if err != nil {
		return nil, err
	}

	//log.Debug("uploadedFiles", uploadedFiles)

	// 落库
	var entFiles []*ent.UserFile
	for _, x := range uploadedFiles {
		entFiles = append(entFiles, &ent.UserFile{
			UserID:   userID,
			Name:     x.Filename,
			Md5:      x.Md5,
			Size:     int64(len(x.Body)),
			Category: filed.FindSuffix(x.Filename),
			//CreatedAt: time.Time{},
			//Extra:     nil,
		})

	}

	err = t.data.Repos.UserFile.InsertBulk(ctx, entFiles)
	if err != nil {
		return nil, err
	}

	// 返回
	var resultFiles []*typepb.File
	for _, x := range uploadedFiles {
		resultFiles = append(resultFiles, &typepb.File{
			//Body: x.Body,
			Name: x.Filename,
			Key:  x.Md5,
			Url:  x.URL,
			Md5:  x.Md5,
		})
	}

	return &databankpb.AddFilesResult{
		Files: resultFiles,
	}, nil
}

func (t *UserFileBiz) asPbFiles(dataFiles data.Files) []typepb.File {
	var files []typepb.File

	for _, x := range dataFiles {
		files = append(files, typepb.File{
			Body: x.Body,
			Name: x.Name,
			Md5:  x.Md5,
			Key:  x.Md5,
			Url:  x.URL,
		})
	}

	return files
}
