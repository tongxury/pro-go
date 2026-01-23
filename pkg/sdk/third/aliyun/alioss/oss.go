package alioss

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-resty/resty/v2"
	"io"
	"store/pkg/sdk/helper"
)

type Config struct {
	AccessKey    string
	AccessSecret string
	Bucket       string
	Endpoint     string
	Region       string
}

func NewClient(conf Config) *Client {
	client, err := oss.New(conf.Endpoint, conf.AccessKey, conf.AccessSecret)
	if err != nil {
		panic(err)
	}
	return &Client{client: client, conf: conf}
}

type Client struct {
	client *oss.Client
	conf   Config
}

func (t *Client) GetSignedUrl(ctx context.Context, bucketName, fileKey, contentType string) (string, error) {

	bucket, err := t.client.Bucket(bucketName)
	if err != nil {
		return "", err
	}

	options := []oss.Option{
		//oss.ObjectACL(oss.ACLPublicRead),
		//oss.Meta("myprop", "mypropval"),
		//oss.ContentType(contentType),
		//ResponseContentEncoding("deflate"),
	}
	url, err := bucket.SignURL(fileKey, oss.HTTPPut, 3600, options...)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (t *Client) UploadByUrl(ctx context.Context, bucketName, url string) (string, error) {

	result, err := resty.New().R().Get(url)
	if err != nil {
		return "", err
	}

	fileBody := bytes.NewReader(result.Body())

	fileBytes, err := io.ReadAll(fileBody)
	if err != nil {
		return "", err
	}

	md5 := helper.MD5(fileBytes)
	//category := filed.FindSuffix(url)

	filename := md5 + ".jpg"

	bucket, err := t.client.Bucket(bucketName)
	if err != nil {
		return "", err
	}

	err = bucket.PutObject(filename, bytes.NewReader(fileBytes))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s.%s/%s", bucketName, t.conf.Endpoint, filename), nil
}

func (t *Client) Upload(ctx context.Context, bucketName, filename string, fileBytes []byte) (string, error) {
	bucket, err := t.client.Bucket(bucketName)
	if err != nil {
		return "", err
	}

	err = bucket.PutObject(filename, bytes.NewReader(fileBytes))
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.%s/%s", bucketName, t.conf.Endpoint, filename)

	return url, nil
}

func (t *Client) UploadBytes(ctx context.Context, filename string, fileBytes []byte) (string, error) {
	bucket, err := t.client.Bucket(t.conf.Bucket)
	if err != nil {
		return "", err
	}

	err = bucket.PutObject(filename, bytes.NewReader(fileBytes))
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.%s/%s", t.conf.Bucket, t.conf.Endpoint, filename)

	return url, nil
}

//https://yoozy.oss-cn-hangzhou.aliyuncs.com/003609743233e391bd8465ec82afc862.jpg
