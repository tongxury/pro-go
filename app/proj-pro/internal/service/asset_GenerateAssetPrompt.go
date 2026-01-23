package service

//
//func (t ProjService) GenerateAssetPrompt() {
//
//	ctx := context.Background()
//
//	list, err := t.data.Mongo.Asset.List(ctx, bson.M{"status": "promptGenerating"})
//	if err != nil {
//		log.Errorw("List err", err)
//		return
//	}
//
//	if len(list) == 0 {
//		return
//	}
//
//	settings, err := t.data.Mongo.Settings.FindOne(ctx, bson.M{})
//	if err != nil {
//		log.Errorw("Settings.FindOne err", err)
//		return
//	}
//
//	data := helper.Mapping(list, func(x *projpb.Asset) *assetPromptJobData {
//		return &assetPromptJobData{
//			Item:     x,
//			Settings: settings,
//		}
//	})
//
//	wg.WaitGroup(ctx, data, t.generateAssetPrompt)
//}
//
//type assetPromptJobData struct {
//	Item     *projpb.Asset
//	Settings *projpb.AppSettings
//}
//
//func (t ProjService) generateAssetPrompt(ctx context.Context, data *assetPromptJobData) error {
//
//	logger := log.NewHelper(log.With(log.DefaultLogger,
//		"func", "generateAssetPrompt",
//		"item", data.Item.XId,
//	))
//
//	prompt, err := t.generateAssetPromptByGemini(ctx, data)
//	if err != nil {
//		logger.Errorw("generateAssetPromptByGemini err", err)
//		prompt, err = t.generateAssetPromptBySeed(ctx, data)
//		if err != nil {
//			logger.Errorw("generateAssetPromptBySeed err", err)
//			return err
//		}
//	}
//
//	_, err = t.data.Mongo.Asset.UpdateByIDIfExists(ctx,
//		data.Item.XId,
//		mgz.Op().
//			Sets(bson.M{
//				"status":                  "promptGenerated",
//				"prompt":                  prompt,
//				"extra.promptGeneratedAt": time.Now().Unix(),
//			}),
//	)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (t ProjService) generateAssetPromptBySeed(ctx context.Context, data *assetPromptJobData) (string, error) {
//
//	asset := data.Item
//	settings := data.Settings.GetVideoScript()
//
//	logger := log.NewHelper(log.With(log.DefaultLogger,
//		"func", "generateAssetPromptBySeed",
//		"item", data.Item.XId,
//	))
//
//	logger.Debugw("start", "")
//
//	segment, err := videoz.GetSegmentByUrl(ctx, asset.Segment.Root.Url, asset.Segment.TimeStart, asset.Segment.TimeEnd)
//	if err != nil {
//		logger.Errorw("GetSegmentByUrl err", err, "start", asset.Segment.TimeStart, "end", asset.Segment.TimeEnd)
//		return "", err
//	}
//
//	base64Video := base64.StdEncoding.EncodeToString(segment.Content)
//
//	var messages []*model.ChatCompletionMessage
//
//	messages = append(messages,
//		&model.ChatCompletionMessage{
//			Role: model.ChatMessageRoleSystem,
//			Content: &model.ChatCompletionMessageContent{
//				StringValue: volcengine.String(settings.Content),
//			},
//		},
//		&model.ChatCompletionMessage{
//			Role: model.ChatMessageRoleSystem,
//			Content: &model.ChatCompletionMessageContent{
//				ListValue: []*model.ChatCompletionMessageContentPart{
//					{
//						Type: "video_url",
//						VideoURL: &model.ChatMessageVideoURL{
//							URL: fmt.Sprintf("data:video/mp4;base64,%s", base64Video),
//						},
//					},
//					{
//						Type: "image_url",
//						ImageURL: &model.ChatMessageImageURL{
//							URL: asset.Commodity.Medias[0].Url,
//						},
//					},
//				},
//			},
//		},
//	)
//
//	req := model.CreateChatCompletionRequest{
//		Model:    "doubao-seed-1-6-250615",
//		Messages: messages,
//	}
//
//	resp, err := t.data.Arkr.C().CreateChatCompletion(ctx, req)
//	if err != nil {
//		logger.Errorw("CreateChatCompletion err", err, "req", req)
//		return "", err
//	}
//
//	prompt := *resp.Choices[0].Message.Content.StringValue
//
//	if asset.Group.GetBaseAssetId() != "" {
//		prompt += fmt.Sprintf(`
//
//## **请一定要继承传给你的视频中的人物和商品**
//`)
//	}
//
//	if asset.Prompts.GetSubtitle() != "" {
//		prompt += fmt.Sprintf(`
//视频中的人物的口播文案一定要用这个：%s
//`, asset.Prompts.GetSubtitle())
//	}
//
//	return prompt, nil
//}
//
//func (t ProjService) generateAssetPromptByGemini(ctx context.Context, data *assetPromptJobData) (string, error) {
//
//	asset := data.Item
//	settings := data.Settings.GetVideoScript()
//
//	log.Debugw("generate asset prompt", "", "asset", asset.XId)
//
//	genaiClient := t.data.GenaiFactory.Get()
//
//	segment, err := videoz.GetSegmentByUrl(ctx, asset.Segment.Root.Url, asset.Segment.TimeStart, asset.Segment.TimeEnd)
//	if err != nil {
//		log.Errorw("GetSegmentByUrl err", err, "start", asset.Segment.TimeStart, "end", asset.Segment.TimeEnd)
//		return "", err
//	}
//
//	genaiUrl, err := genaiClient.UploadBlob(ctx, segment.Content, "video/mp4")
//	if err != nil {
//		return "", err
//	}
//
//	image := asset.Commodity.Medias[0].Url
//
//	genaiUrlImage, err := genaiClient.UploadFile(ctx, image, "image/jpeg")
//	if err != nil {
//		return "", err
//	}
//
//	parts := []*genai.Part{
//		{
//			FileData: &genai.FileData{
//				MIMEType: "video/mp4",
//				FileURI:  genaiUrl,
//			},
//		},
//		{
//			FileData: &genai.FileData{
//				MIMEType: "image/jpeg",
//				FileURI:  genaiUrlImage,
//			},
//		},
//
//		{
//			Text: fmt.Sprintf(settings.Content),
//		},
//	}
//
//	prompt, err := genaiClient.GenerateContent(ctx, gemini.GenerateContentRequest{
//		Model: gemini.DefaultModel,
//		Parts: parts,
//	})
//
//	log.Debugw("prompt", prompt)
//
//	if err != nil {
//		return "", err
//	}
//
//	if prompt == "" {
//		log.Errorw("prompt err", "empty prompt")
//		return "", nil
//	}
//
//	if asset.Group.GetBaseAssetId() != "" {
//		prompt += fmt.Sprintf(`
//
//## **请一定要继承传给你的视频中的人物和商品**
//`)
//	}
//
//	if asset.Prompts.GetSubtitle() != "" {
//		prompt += fmt.Sprintf(`
//视频中的人物的口播文案一定要用这个：%s
//`, asset.Prompts.GetSubtitle())
//	}
//
//	return prompt, nil
//}
