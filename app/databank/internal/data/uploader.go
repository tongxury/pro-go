package data

//type Uploader struct {
//	s3Uploader *s3manager.Uploader
//}
//
//func NewUploader(s3Uploader *s3manager.Uploader) *Uploader {
//	return &Uploader{
//		s3Uploader: s3Uploader,
//	}
//}
//
//func (t Uploader) Upload(ctx context.Context, files kratosutil.FormFiles, dir string) (kratosutil.FormFiles, error) {
//
//	getPath := func(dir string, md5, name string) string {
//		return dir + "/" + md5 + filed.FindSuffix(name)
//	}
//
//	bucket := "studygpt-pub"
//
//	objects := make([]s3manager.BatchUploadObject, 0, len(files))
//
//	for _, x := range files {
//
//		//path := helper.MD5(x.Body) + filed.FindSuffix(x.Name)
//
//		md5 := helper.MD5(x.Body)
//		path := getPath(dir, md5, x.Filename)
//
//		x.URL = fmt.Sprintf("https://%s.s3.amazonaws.com/", bucket) + getPath(dir, md5, x.Filename)
//		x.Md5 = md5
//
//		objects = append(objects, s3manager.BatchUploadObject{
//			Object: &s3manager.UploadInput{
//				Field:    aws.String(path),
//				Bucket: aws.String(bucket),
//				Body:   bytes.NewReader(x.Body),
//				ACL:    aws.String(storagegateway.ObjectACLPublicRead),
//			},
//		})
//	}
//
//	iter := &s3manager.UploadObjectsIterator{Objects: objects}
//	err := t.s3Uploader.UploadWithIterator(ctx, iter)
//	if err != nil {
//		return nil, err
//	}
//
//	return files, nil
//}
