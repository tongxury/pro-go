package data

import (
	"context"
	"store/pkg/krathelper"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/filed"
	"store/pkg/sdk/third/aliyun/alioss"
)

type Uploader struct {
	client *alioss.Client
}

func NewUploader(client *alioss.Client) *Uploader {
	return &Uploader{
		client: client,
	}
}

func (t Uploader) Upload(ctx context.Context, files krathelper.FormFiles, dir string) (krathelper.FormFiles, error) {

	getPath := func(dir string, md5, name string) string {
		return dir + "/" + md5 + filed.FindSuffix(name)
	}

	var results krathelper.FormFiles
	for _, x := range files {

		md5 := helper.MD5(x.Body)
		objectKey := getPath(dir, md5, x.Filename)

		url, err := t.client.Upload(ctx, "xtips", objectKey, x.Body)
		if err != nil {
			return nil, err
		}

		x.URL = url
		x.Md5 = md5
		results = append(results, x)
	}

	return results, nil
}
