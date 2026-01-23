package tos

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"net/http"
	"store/pkg/sdk/helper"
	"time"

	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
)

type Config struct {
	AccessKeyID     string
	AccessKeySecret string
	Endpoint        string
	Region          string
	DefaultBucket   string
}
type Client struct {
	c    *tos.ClientV2
	conf Config
}

func NewClient(conf Config) *Client {

	credential := tos.NewStaticCredentials(conf.AccessKeyID, conf.AccessKeySecret)
	// 可以通过 tos.WithXXX 的方式添加可选参数
	// 示例中通过 WithConnectionTimeout 设置建立连接超时时间
	// 通过 WithSocketTimeout 设置一次读写连接超时时间
	client, err := tos.NewClientV2(conf.Endpoint,
		tos.WithCredentials(credential),
		tos.WithRegion(conf.Region),
		tos.WithConnectionTimeout(300*time.Second),
		tos.WithSocketTimeout(300*time.Second, 300*time.Second))
	if err != nil {
		panic(err)
	}

	// 使用结束后，关闭 client
	//client.Close()
	return &Client{
		c:    client,
		conf: conf,
	}
}

type PutByUrlRequest struct {
	Bucket string
	Url    string
	Suffix string
}

func (t *Client) PutVideo(ctx context.Context, url string) (string, error) {

	return t.PutByUrl(ctx, PutByUrlRequest{
		Bucket: t.conf.DefaultBucket,
		Url:    url,
		Suffix: ".mp4",
	})
}

func (t *Client) PutImage(ctx context.Context, url string) (string, error) {

	return t.PutByUrl(ctx, PutByUrlRequest{
		Bucket: t.conf.DefaultBucket,
		Url:    url,
		Suffix: ".jpg",
	})
}

func (t *Client) PutImageBytes(ctx context.Context, bytes []byte) (string, error) {

	return t.Put(ctx, PutRequest{
		Bucket:  t.conf.DefaultBucket,
		Content: bytes,
		Key:     helper.MD5(bytes) + ".jpg",
	})
}

func (t *Client) PutVideoBytes(ctx context.Context, bytes []byte) (string, error) {

	return t.Put(ctx, PutRequest{
		Bucket:  t.conf.DefaultBucket,
		Content: bytes,
		Key:     helper.MD5(bytes) + ".mp4",
	})
}

func (t *Client) PutAudioBytes(ctx context.Context, bytes []byte) (string, error) {

	return t.Put(ctx, PutRequest{
		Bucket:  t.conf.DefaultBucket,
		Content: bytes,
		Key:     helper.MD5(bytes) + ".mp3",
	})
}

func (t *Client) PutByUrl(ctx context.Context, req PutByUrlRequest) (string, error) {

	v, err := http.Get(req.Url)
	if err != nil {
		return "", nil
	}

	all, err := io.ReadAll(v.Body)
	if err != nil {
		return "", err
	}

	if len(all) == 0 {
		return "", errors.New("empty body")
	}

	md5Value := fmt.Sprintf("%x", md5.Sum(all))

	key := md5Value + req.Suffix

	_, err = t.c.PutObjectV2(ctx, &tos.PutObjectV2Input{
		PutObjectBasicInput: tos.PutObjectBasicInput{
			Bucket:        req.Bucket,
			ContentLength: int64(len(all)),
			Key:           key,
			//ContentMD5: md5Value,
		},
		Content: bytes.NewReader(all),
	})
	if err != nil {
		return "", err
	}

	//https: //yoozyres.tos-cn-shanghai.volces.com/a18afd25c51b10d50e5c1dd581de027c.jpg

	return fmt.Sprintf("https://%s.%s/%s", req.Bucket, t.conf.Endpoint, key), nil
}

type PutRequest struct {
	Bucket  string
	Content []byte
	Key     string
}

func (t *Client) Put(ctx context.Context, req PutRequest) (string, error) {

	bucket := helper.OrString(req.Bucket, t.conf.DefaultBucket)

	_, err := t.c.PutObjectV2(ctx, &tos.PutObjectV2Input{
		PutObjectBasicInput: tos.PutObjectBasicInput{
			Bucket: bucket,
			Key:    req.Key,
		},
		Content: bytes.NewReader(req.Content),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s.%s/%s", bucket, t.conf.Endpoint, req.Key), nil

}
