package data

import (
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

func (t Uploader) Client() *alioss.Client {
	return t.client
}
