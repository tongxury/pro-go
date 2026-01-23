package gemini

//func (t *Client) UploadBlob(ctx context.Context, content []byte, mimeType string) (string, error) {
//
//	// 多个client间不能共享 file链接
//	//cacheKey := "gemini:" + helper.MD5(content)
//	//get, err := t.cache.Get(ctx, cacheKey)
//	//if err == nil && get != "" {
//	//	return get, err
//	//}
//	//imgData, _ := os.ReadFile(imagePath)
//
//	now := time.Now()
//
//	fileBody := bytes.NewReader(content)
//
//	ctx = context.Background()
//
//	opts := genai.UploadFileConfig{DisplayName: "", MIMEType: mimeType}
//
//	var response *genai.File
//	var err error
//
//	for i := 0; i < 5; i++ {
//		response, err = t.c.Files.Upload(ctx, fileBody, &opts)
//		if err == nil {
//			break
//		}
//
//		sleep := time.Duration(1<<i) * time.Second
//		log.Errorw("upload file error", err, "retry", i, "sleep", sleep)
//		time.Sleep(sleep)
//		fileBody.Seek(0, 0) // Reset reader
//	}
//
//	if err != nil {
//		log.Errorw("upload file failed after retries", err)
//		return "", err
//	}
//
//	time.Sleep(5 * time.Second)
//
//	for {
//		time.Sleep(1 * time.Second)
//
//		fmt.Println("checking", response.Name, response.State)
//
//		response, err = t.c.Files.Get(ctx, response.Name, nil)
//		if err != nil {
//			log.Fatal(err)
//			continue
//		}
//
//		if response.State == genai.FileStateProcessing {
//			continue
//		}
//
//		if response.State == genai.FileStateActive {
//
//			//t.cache.Set(ctx, cacheKey, response.URI, 5*time.Hour)
//
//			fmt.Println(time.Now().Sub(now).Milliseconds(), len(content))
//
//			return response.URI, nil
//		}
//
//		if response.State == genai.FileStateFailed {
//			return "", errors.New(response.Error.Message)
//		}
//	}
//}
//
//func (t *Client) UploadFile(ctx context.Context, url, mimeType string) (string, error) {
//
//	client := resty.New().SetTLSClientConfig(&tls.Config{
//		InsecureSkipVerify: true,
//	})
//
//	result, err := client.R().Get(url)
//	if err != nil {
//		return "", err
//	}
//
//	return t.UploadBlob(ctx, result.Body(), mimeType)
//}
