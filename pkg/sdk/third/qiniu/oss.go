package qiniu

import (
	"context"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
)

type Config struct {
	AccessKey    string
	AccessSecret string
	//Bucket       string
	//Endpoint     string
}

func NewClient(conf Config) *Client {
	return &Client{conf: conf}
}

type Client struct {
	conf Config
}

func (t Client) GetUploadTokenV2(ctx context.Context, bucket string) (string, error) {

	mac := credentials.NewCredentials(t.conf.AccessKey, t.conf.AccessSecret)

	putPolicy, err := uptoken.NewPutPolicy(bucket, time.Now().Add(1*time.Hour))
	if err != nil {
		return "", err
	}
	upToken, err := uptoken.NewSigner(putPolicy, mac).GetUpToken(ctx)
	if err != nil {
		return "", err
	}

	return upToken, nil
}

func (t Client) GetUploadToken(ctx context.Context, bucket string) (string, error) {

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := auth.New(t.conf.AccessKey, t.conf.AccessSecret)
	upToken := putPolicy.UploadToken(mac)

	return upToken, nil
}
