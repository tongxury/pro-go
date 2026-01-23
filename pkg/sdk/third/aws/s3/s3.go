package s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/storagegateway"
)

type Client struct {
	conf Config
	*s3manager.Uploader
}

func NewS3Client(conf Config) *Client {

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(conf.AccessKey, conf.AccessSecret, ""),
		//Endpoint:         aws.String(end_point),
		Region:           aws.String(conf.Region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(false),
	})

	if err != nil {
		panic(err)
	}

	svc := s3.New(sess)

	uploader := s3manager.NewUploaderWithClient(svc, func(u *s3manager.Uploader) {
		// 定义将在内存中缓冲25个MiB的策略
		//u.BufferProvider = s3manager.NewBufferedReadSeekerWriteToPool(25 * 1024 * 1024)
		//u.PartSize = 64 * 1024 * 1024 // 每个部分 64MB
	})
	//objects := make([]s3manager.BatchUploadObject, 0, len(keys))
	//for _, key := range keys {
	//	fOpen, err := os.Open(key)
	//	if err != nil {
	//		panic(err)
	//	}
	//	defer fOpen.Close()
	//	objects = append(objects, s3manager.BatchUploadObject{
	//		Object: &s3manager.UploadInput{
	//			Body:   fOpen,
	//			Bucket: aws.String(bucket),
	//			Field:    aws.String(key),
	//		},
	//	})
	//}
	//
	//bytes.NewReader()
	//
	//iter := &s3manager.UploadObjectsIterator{Objects: objects}
	//if err := uploader.UploadWithIterator(aws.BackgroundContext(), iter); err != nil {
	//	panic(err)
	//}
	return &Client{
		conf:     conf,
		Uploader: uploader,
	}
}

// "oscar-res", md5+category, fileBytes
func (t *Client) Upload(ctx context.Context, mimeType, filename string, fileBytes []byte) (string, error) {

	objects := make([]s3manager.BatchUploadObject, 0, 1)

	//https://veogo-resources.s3.us-east-1.amazonaws.com/6c22ea32abc79477f19c3da71488a636.jpg

	objects = append(objects, s3manager.BatchUploadObject{
		Object: &s3manager.UploadInput{
			//Metadata: map[string]*string{
			//	"Content-Type": aws.String(filed.ContentTypeByName(x.Filename)),
			//},
			ContentType: aws.String(mimeType),
			//ContentType: aws.String("image/png"),
			Key:    aws.String(filename),
			Bucket: aws.String(t.conf.Bucket),
			Body:   bytes.NewReader(fileBytes),
			ACL:    aws.String(storagegateway.ObjectACLPublicRead),
			//ContentMD5:  aws.String(md5),
		},
	})

	iter := &s3manager.UploadObjectsIterator{Objects: objects}
	err := t.UploadWithIterator(ctx, iter)
	if err != nil {
		return "", err
	}

	//return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", t.conf.Bucket, t.conf.Region, filename), nil
	return fmt.Sprintf("%s/%s", t.conf.Endpoint, filename), nil
}
