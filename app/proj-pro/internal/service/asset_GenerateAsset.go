package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	creditpb "store/api/credit"
	projpb "store/api/proj"
	"store/app/proj-pro/configs"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/imagez"
	"store/pkg/sdk/helper/videoz"
	"store/pkg/sdk/helper/wg"
	"store/pkg/sdk/third/bytedance/tos"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/openai/openai-go/v3"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjService) GenerateAsset() {

	ctx := context.Background()

	list, err := t.data.Mongo.Asset.List(ctx, bson.M{"status": "generating"})
	if err != nil {
		log.Errorw("List err", err)
		return
	}

	if len(list) == 0 {
		return
	}

	settings, err := t.data.Mongo.Settings.FindOne(ctx, bson.M{})
	if err != nil {
		log.Errorw("Settings.FindOne err", err)
		return
	}

	data := helper.Mapping(list, func(x *projpb.Asset) *assetJobData {
		return &assetJobData{
			Item:     x,
			Settings: settings,
		}
	})

	wg.WaitGroup(ctx, data, t.generateAssetVideo)
}

type assetJobData struct {
	Item     *projpb.Asset
	Settings *projpb.AppSettings
}

func (t ProjService) generateAssetVideo(ctx context.Context, data *assetJobData) error {
	//if data.Settings.GetVideoGenerate().GetModel() == "seed" {
	//return t.generateAssetVideoBySeed(ctx, data)
	//}

	return t.generateAssetVideoBySora2(ctx, data)
}

func (t ProjService) generateAssetVideoBySeed(ctx context.Context, data *assetJobData) error {

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "generateAssetVideoBySeed",
		"item", data.Item.XId,
	))

	asset := data.Item

	logger.Debug("start")

	processingSeedVideoId := helper.FindInStringMap(asset.Extra.GetContext(), "processingSeedVideoId")
	if processingSeedVideoId != "" {
		task, err := t.data.Arkr.C().GetContentGenerationTask(ctx, model.GetContentGenerationTaskRequest{
			ID: processingSeedVideoId,
		})
		if err != nil {
			logger.Errorw("GetContentGenerationTask err", err, "processingSeedVideoId", processingSeedVideoId)
			return err
		}

		if task.Status == "failed" {

			failedCount := conv.Int(helper.FindInStringMap(asset.Extra.GetContext(), "failedCount")) + 1

			if failedCount > 2 {
				t.data.Mongo.Asset.UpdateByIDXX(ctx,
					asset.XId,
					mgz.Set(bson.M{
						"extra.context.processingSeedVideoId": "",
						"extra.context.failedReason":          task.Error.Message,
						"status":                              "failed",
						"extra.context.status":                task.Status,
						"extra.model":                         data.Settings.GetVideoGenerate().GetModel(),
					}))

				return nil
			}

			t.data.Mongo.Asset.UpdateByIDXX(ctx,
				asset.XId,
				mgz.Set(bson.M{
					"extra.context.processingSeedVideoId": "",
					"extra.context.failedReason":          task.Error.Message,
					"extra.context.failedCount":           conv.Str(failedCount),
				}))

			return nil
		}

		if task.Status == "succeeded" {

			url, err := t.data.TOS.PutByUrl(ctx, tos.PutByUrlRequest{
				Bucket: "yoozyres",
				Url:    task.Content.VideoURL,
				Suffix: ".mp4",
			})
			if err != nil {
				return err
			}

			coverUrl, err := t.data.TOS.PutByUrl(ctx, tos.PutByUrlRequest{
				Bucket: "yoozyres",
				Url:    task.Content.LastFrameURL,
				Suffix: ".jpg",
			})
			if err != nil {
				return err
			}

			_, err = t.data.Mongo.Asset.UpdateByIDXX(ctx,
				asset.XId,
				mgz.Set(bson.M{
					"coverUrl": coverUrl,
					"url":      url,
					"status":   "completed",
					//"duration":                        conv.Float64(task.Seed),
					"extra.context.processingSeedVideoId": "",
					"extra.context.seedVideoId":           processingSeedVideoId,
					"extra.completedAt":                   time.Now().Unix(),
					"extra.model":                         "seed",
				}))
			if err != nil {
				log.Errorw("UpdateByIDXX err", err)
				return err
			}

			return nil
		}

		// 其他情况 继续轮询
		_, err = t.data.Mongo.Asset.UpdateByIDXX(ctx,
			asset.XId,
			mgz.Set(bson.M{
				"extra.context.status": task.Status,
			}))

		if err != nil {
			return err
		}

		time.Sleep(30 * time.Second)

		return nil
	}

	// 普通生成视频

	if asset.Prompt == "" {
		_ = t.data.Mongo.Asset.DeleteByID(ctx, asset.XId)
		return nil
	}

	req := model.CreateContentGenerationTaskRequest{
		Model: "doubao-seedance-1-0-pro-250528", // Replace with Model ID
		Content: []*model.CreateContentGenerationContentItem{
			{
				Type: "text",
				Text: volcengine.String(asset.Prompt + "--rt 9:16 --rs 720p --dur 10 --cf false"),
			},
			{
				Type: "image_url",
				ImageURL: &model.ImageURL{
					URL: asset.Commodity.Medias[0].Url,
				},
			},
		},
	}

	//content := []*model.CreateContentGenerationContentItem{
	//	{
	//		Type: "text",
	//		Text: volcengine.String(x + "--rs 720p --dur 10 --cf false"),
	//	},
	//}
	//for _, xx := range helper.SubSlice(task.Commodity.GetImages(), 3) {
	//	content = append(content, &model.CreateContentGenerationContentItem{
	//		Type: "image_url",
	//		ImageURL: &model.ImageURL{
	//			URL: xx,
	//		},
	//		Role: volcengine.String("reference_image"),
	//	})
	//}
	//
	//req := model.CreateContentGenerationTaskRequest{
	//	//Model: "doubao-seedance-1-0-pro-250528",
	//	Model:   "doubao-seedance-1-0-lite-i2v-250428",
	//	Content: content,
	//}
	//

	resp, err := t.data.Arkr.C().CreateContentGenerationTask(ctx, req)
	if err != nil {
		logger.Errorw("CreateContentGenerationTask err", err)
		return nil
	}

	t.data.Mongo.Asset.UpdateByIDXX(ctx,
		asset.XId,
		mgz.Set(bson.M{
			"extra.context.processingSeedVideoId": resp.ID,
		}),
	)

	return nil
}

func (t ProjService) generateAssetVideoBySora2(ctx context.Context, data *assetJobData) error {

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "generateAssetVideoBySora2",
		"item", data.Item.XId,
	))

	asset := data.Item

	logger.Debug("start")

	//log.Debugw("generateAssetVideo", "ing", "asset", asset.XId)

	processingVideoId := helper.FindInStringMap(asset.Extra.GetContext(), "processingVideoId")
	if processingVideoId != "" {
		status, err := t.data.OpenAI.Videos().Get(ctx, processingVideoId)
		if err != nil {
			logger.Errorw("OpenAI.PollStatus err", err, "processingVideoId", processingVideoId)
			return err
		}

		logger.Debugw("PollStatus ", "", "videoId", status.ID, "status", status.Status, "process", status.Progress)

		// 失败了重新生成视频
		if status.Status == openai.VideoStatusFailed {

			failedCount := conv.Int(helper.FindInStringMap(asset.Extra.GetContext(), "failedCount")) + 1

			if failedCount > 2 {
				t.data.Mongo.Asset.UpdateByIDXX(ctx, asset.XId,
					mgz.Set(bson.M{
						"extra.context.processingVideoId": "",
						"extra.context.failedReason":      status.Error.Message,
						"status":                          "failed",
						"extra.context.status":            status.Status,
						"extra.model":                     data.Settings.GetVideoGenerate().GetModel(),
					}))

				return nil
			}

			t.data.Mongo.Asset.UpdateByIDXX(ctx, asset.XId,
				mgz.Set(bson.M{
					"extra.context.processingVideoId": "",
					"extra.context.failedReason":      status.Error.Message,
					"extra.context.failedCount":       conv.Str(failedCount),
				}))

			return nil
		}

		// 成功了结束任务
		if status.Status == openai.VideoStatusCompleted {
			cr, err := t.data.OpenAI.Videos().DownloadContent(ctx, processingVideoId, openai.VideoDownloadContentParams{
				Variant: "video",
			})
			if err != nil {
				return err
			}

			content, err := io.ReadAll(cr.Body)
			if err != nil {
				return err
			}

			content, err = videoz.RemoveFirstFrame(content)
			if err != nil {
				return err
			}
			//url, err := t.data.Alioss.UploadBytes(ctx, helper.MD5(content)+".mp4", content)
			url, err := t.data.TOS.Put(ctx, tos.PutRequest{
				Bucket:  "yoozyres",
				Content: content,
				Key:     helper.MD5(content) + ".mp4",
			})
			if err != nil {
				return err
			}

			logger.Debugw("Upload video ", url)

			coverFrame, err := videoz.GetFrame(content, 1)
			if err != nil {
				log.Errorw("GetFrame err", err)
				return err
			}

			logger.Debugw("GetCoverFrame url", url)

			coverUrl, err := t.data.TOS.Put(ctx, tos.PutRequest{
				Bucket:  "yoozyres",
				Content: coverFrame,
				Key:     helper.MD5(content) + ".jpg",
			})
			if err != nil {
				log.Errorw("Put err", err)
				return err
			}

			_, err = t.data.Mongo.Asset.UpdateByIDXX(ctx,
				asset.XId,
				mgz.Set(bson.M{
					"coverUrl":                        coverUrl,
					"url":                             url,
					"status":                          "completed",
					"duration":                        conv.Float64(status.Seconds),
					"extra.context.processingVideoId": "",
					"extra.context.videoId":           processingVideoId,
					"extra.completedAt":               time.Now().Unix(),
					"extra.model":                     "sora2",
				}))
			if err != nil {
				log.Errorw("UpdateByIDXX err", err)
				return err
			}

			logger.Debugw("Update status ", "success")

			_, err = t.data.GrpcClients.CreditClient.XCost(ctx, &creditpb.XCostRequest{
				UserId: asset.GetUser().GetXId(),
				Amount: configs.CreditCostAsset,
				Key:    asset.XId,
			})
			if err != nil {
				logger.Errorw("Cost err", err, "userId", asset.GetUser().GetXId())
			}

			logger.Debugw("XCost", "", "userId", asset.GetUser().GetXId(), "amount", configs.CreditCostAsset)

			return nil
		}

		// 其他情况 继续轮询
		// 配合前端展示
		statusT := status.Status
		if statusT == "in_progress" {
			statusT = "running"
		}

		_, err = t.data.Mongo.Asset.UpdateByIDXX(ctx, asset.XId, mgz.Set(bson.M{
			"extra.context.status": statusT,
		}))
		if err != nil {
			return err
		}

		time.Sleep(30 * time.Second)

		return nil
	}

	// 视频调整
	if asset.Group.GetRefAssetId() != "" {

		refAsset, err := t.data.Mongo.Asset.GetById(ctx, asset.Group.GetRefAssetId())
		if err != nil {
			return err
		}

		videoId := helper.FindInStringMap(refAsset.Extra.GetContext(), "videoId")
		if videoId == "" {
			return errors.New("empty videoId")
		}

		// 调整视频
		log.Debugw("generateAssetVideo ", "Remix", "asset", asset.XId, "videoId", videoId, "prompt", asset.PromptAddition)

		video, err := t.data.OpenAI.Videos().Remix(ctx, videoId, openai.VideoRemixParams{
			Prompt: helper.OrString(asset.Prompts.GetScriptAdjustment(), asset.PromptAddition),
		})

		if err != nil {
			log.Errorw("start generate error", err, "asset", asset)

			if video != nil && video.Error.Code == "rate_limit_exceeded" {
				time.Sleep(1 * time.Minute)
			}

			return err
		}

		t.data.Mongo.Asset.UpdateByIDXX(ctx, asset.XId, mgz.Set(bson.M{
			"extra.context.processingVideoId": video.ID,
		}))

		return nil
	}

	prompt := helper.OrString(asset.Prompts.GetScript(), asset.Prompt)

	// base视频生成
	if asset.Group.GetBaseAssetId() != "" {

		baseAsset, err := t.data.Mongo.Asset.GetById(ctx, asset.Group.GetRefAssetId())
		if err != nil {
			return err
		}

		videoId := helper.FindInStringMap(baseAsset.Extra.GetContext(), "videoId")
		if videoId == "" {
			return errors.New("empty videoId")
		}

		// 调整视频
		log.Debugw("generateAssetVideo ", "Remix", "asset", asset.XId, "videoId", videoId, "prompt", asset.PromptAddition)

		video, err := t.data.OpenAI.Videos().Remix(ctx, videoId, openai.VideoRemixParams{
			Prompt: prompt,
		})

		if err != nil {
			log.Errorw("start generate error", err, "asset", asset)

			if video != nil && video.Error.Code == "rate_limit_exceeded" {
				time.Sleep(1 * time.Minute)
			}

			return err
		}

		t.data.Mongo.Asset.UpdateByIDXX(ctx, asset.XId, mgz.Set(bson.M{
			"extra.context.processingVideoId": video.ID,
		}))

		return nil
	}

	// 普通生成视频

	if asset.Prompt == "" {
		_ = t.data.Mongo.Asset.DeleteByID(ctx, asset.XId)
		return nil
	}

	seconds := "8"
	if asset.Segment.TimeEnd-asset.Segment.TimeStart > 8 {
		seconds = "12"
	}

	imgBytes, err := helper.ReadBytesByUrl(asset.Commodity.Medias[0].Url)
	if err != nil {
		return err
	}

	resizedImage, err := imagez.ResizeKeepRatio(imgBytes, 720, 1280)
	if err != nil {
		return err
	}

	video, err := t.data.OpenAI.Videos().New(ctx, openai.VideoNewParams{
		Model:          openai.VideoModelSora2,
		Prompt:         prompt,
		InputReference: openai.File(bytes.NewReader(resizedImage), helper.CreateUUID(), "image/jpeg"),
		Seconds:        openai.VideoSeconds(seconds),
	})

	if err != nil {
		log.Errorw("start generate error", err, "asset", asset)

		if video != nil && video.Error.Code == "rate_limit_exceeded" {
			time.Sleep(1 * time.Minute)
		}

		return err
	}

	t.data.Mongo.Asset.UpdateByIDXX(ctx, asset.XId, mgz.Set(bson.M{
		"extra.context.processingVideoId": video.ID,
	}))

	return nil
}

//func (t ProjService) generateAsset2(ctx context.Context, asset *projpb.Asset) error {
//
//	log.Debugw("generate asset", "", "asset", asset.XId)
//
//	genaiClient := t.data.GenaiFactory.Get()
//
//	segment, err := videoz.GetSegmentByUrl(ctx, asset.Segment.Root.Url, asset.Segment.TimeStart, asset.Segment.TimeEnd)
//	if err != nil {
//		log.Errorw("GetSegmentByUrl err", err, "start", asset.Segment.TimeStart, "end", asset.Segment.TimeEnd)
//		return err
//	}
//
//	genaiUrl, err := genaiClient.UploadBlob(ctx, segment.Content, "video/mp4")
//	if err != nil {
//		return err
//	}
//
//	image := asset.Commodity.Medias[0].Url
//
//	genaiUrlImage, err := genaiClient.UploadFile(ctx, image, "image/jpeg")
//	if err != nil {
//		return err
//	}
//
//	models := []string{
//		"gemini-2.5-pro",
//		"gemini-2.5-pro-preview-05-06",
//		//"gemini-2.5-flash",
//		//"gemini-2.5-flash-lite-preview-06-17",
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
//			Text: fmt.Sprintf(`
//基于我提供的参考视频，帮我写一个提示词给到Sora2模型来完全复刻这个视频，包括视频中人物讲的口播文案和对应的语言以及视频的环境音或者配乐。这个视频的目标是用来做抖音的短视频的第一个吸引眼球的镜头。手机拍摄视角。注意你只能回复提示词的内容。
//
//注意：
//1.提示词要写成分镜头脚本；
//2.口播文案要判断视频中是说中文还是英文，口播文案需要明确在提示词中，口播文案需要指明是哪个角色说的，以及什么时候说的；
//3.视频时长限制在15秒内；
//4.保持光线与配色一致以便剪辑，不需要展示字幕;
//5.场景需要真实自然；
//6.人物动作需要真实自然。
//7.提示词要详细说明人物妆容、动作、神态，便于复刻原视频
//8.提示词要详细说明画面构图，便于复刻原视频
//9.提示词要详细说明转场及运镜方式，便于复刻原视频
//10.**提示词中的商品要换成给定的参考图中的商品**
//
//`),
//		},
//	}
//
//	prompt, err := genaiClient.GenerateContent(ctx, gemini.GenerateContentRequest{
//		Model: models[mathz.RandNumber(0, len(models)-1)],
//		Parts: parts,
//	})
//
//	if asset.Prompt != "" {
//		prompt += fmt.Sprintf(`
//商品要用给定的参考图中的商品
//
//附加要求: %s
//`, asset.Prompt)
//	}
//
//	log.Debugw("prompt", prompt)
//
//	seconds := "4"
//	if asset.Segment.TimeEnd-asset.Segment.TimeStart > 4 {
//		seconds = "8"
//	}
//
//	tmpImages, err := t.data.Arkr.C().GenerateImages(ctx, model.GenerateImagesRequest{
//		Model: "doubao-seedream-4-0-250828",
//		//Prompt: "将给定的所有图片中涉及到的商品替换成[图一]中的商品，并将替换后的图片按顺序合并展示在一张图中, 去掉人像，不要裁剪",
//		Prompt: `
//帮我裁剪，不要做其他任何动作。
//`,
//		//Prompt:         "输出包含给定图商品的图片",
//		Image:          image,
//		Size:           volcengine.String("720x1280"),
//		ResponseFormat: volcengine.String(model.GenerateImagesResponseFormatURL),
//		Watermark:      volcengine.Bool(false),
//	})
//	if err != nil {
//		log.Errorw("GenerateImages err", "err", err)
//		return err
//	}
//
//	response, err := http.Get(*tmpImages.Data[0].Url)
//	if err != nil {
//		return err
//	}
//
//	video, err := t.data.OpenAI.Videos().NewAndPoll(ctx, openai.VideoNewParams{
//		Model:          openai.VideoModelSora2,
//		Prompt:         prompt,
//		InputReference: openai.File(response.Body, helper.CreateUUID(), "image/jpeg"),
//		Seconds:        openai.VideoSeconds(seconds),
//	}, 1000)
//
//	if err != nil {
//		log.Errorw("start generate error", err, "asset", asset)
//		return err
//	}
//
//	if video.Status == openai.VideoStatusFailed {
//		log.Errorf("Video creation failed. Status: %s\n", video.Status)
//
//		//retryCount := conv.Int(asset.Extra["count"])
//		if asset.Extra.GetRetryCount() > 10 {
//			t.data.Mongo.Asset.UpdateByIDXX(ctx, asset.XId, mgz.Set(bson.M{
//				"status":       "failed",
//				"extra.reason": video.Error,
//			}))
//		}
//
//		t.data.Mongo.Asset.UpdateByIDXX(ctx, asset.XId, mgz.Set(bson.M{
//			"extra.retryCount": asset.Extra.GetRetryCount() + 1,
//			"extra.reason":     video.Error,
//		}))
//
//		return fmt.Errorf("Video creation failed. Status: %s\n", video.Status)
//
//	}
//
//	cr, err := t.data.OpenAI.Videos().DownloadContent(ctx, video.ID, openai.VideoDownloadContentParams{
//		Variant: "video",
//	})
//	if err != nil {
//		return err
//	}
//
//	content, err := io.ReadAll(cr.Body)
//	if err != nil {
//		return err
//	}
//
//	//videoz.GetFrame()
//
//	//url, err := t.data.Alioss.UploadBytes(ctx, helper.MD5(content)+".mp4", content)
//	//if err != nil {
//	//	return err
//	//}
//
//	url, err := t.data.TOS.Put(ctx, tos.PutRequest{
//		Bucket:  "yoozyres",
//		Content: content,
//		Key:     helper.MD5(content) + ".mp4",
//	})
//	if err != nil {
//		return err
//	}
//
//	_, err = t.data.Mongo.Asset.UpdateByIDXX(ctx, asset.XId, mgz.Set(bson.M{
//		"coverUrl": image,
//		"url":      url,
//		"status":   "completed",
//		"duration": conv.Float64(seconds),
//	}))
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func (t ProjService) generateAsset(ctx context.Context, asset *projpb.Asset) error {

	images := []string{asset.Commodity.Medias[0].Url}

	for _, x := range asset.Segment.HighlightFrames {
		images = append(images, x.Url)
	}

	tmpImages, err := t.data.Arkr.C().GenerateImages(ctx, model.GenerateImagesRequest{
		Model: "doubao-seedream-4-0-250828",
		//Prompt: "将给定的所有图片中涉及到的商品替换成[图一]中的商品，并将替换后的图片按顺序合并展示在一张图中, 去掉人像，不要裁剪",
		Prompt: `
所有图片按顺序缩放展示，去掉人像，不要对原图片进行裁剪，只能缩放，不要做融合，只做拼接。
`,
		//Prompt:         "输出包含给定图商品的图片",
		Image:          images,
		Size:           volcengine.String("720x1280"),
		ResponseFormat: volcengine.String(model.GenerateImagesResponseFormatURL),
		Watermark:      volcengine.Bool(false),
	})
	if err != nil {
		return err
	}

	response, err := http.Get(*tmpImages.Data[0].Url)
	if err != nil {
		return err
	}
	prompt := fmt.Sprintf(
		`
	按照给定参考图中体现的分镜序列图生成抖音带货视频。
	
	画面描述:%s
	--
	整体风格：%s
	--
	画面风格：%s
	--
	场景风格：%s
	---
	标签：%s
	===
	整体要求：
	1. 口播文案：[%s]，可适当增加语速，一定要把文案讲完
	2. 不要添加任何字幕
	3. 保持光线一致以便后续剪辑
	4. 场景需要真实自然；
	5.人物动作需要真实自然。
	6. 这是一个抖音带货短视频，节奏需要紧凑
	7. 生成视频中的商品一定要用给定参考图中优先出现的商品
	`,

		//7. 只参考给定图的叙事序列，不要直接将参考图作为生成视频的某一帧

		asset.Segment.Description,
		asset.Segment.Style,
		asset.Segment.ContentStyle,
		asset.Segment.SceneStyle,
		conv.S2J(asset.Segment.TypedTags),
		asset.Segment.Subtitle,
	)

	if asset.Prompt != "" {
		prompt += fmt.Sprintf(`
附加要求: %s
`, asset.Prompt)
	}

	log.Debugw("prompt", prompt)

	seconds := "4"
	if asset.Segment.TimeEnd-asset.Segment.TimeStart > 4 {
		seconds = "8"
	}

	video, err := t.data.OpenAI.Videos().NewAndPoll(ctx, openai.VideoNewParams{
		Model:          openai.VideoModelSora2,
		Prompt:         prompt,
		InputReference: openai.File(response.Body, helper.CreateUUID(), "image/jpeg"),
		Seconds:        openai.VideoSeconds(seconds),
	}, 1000)

	if err != nil {
		log.Errorw("start generate error", err, "asset", asset)
		return err
	}

	if video.Status == openai.VideoStatusFailed {
		log.Errorf("Video creation failed. Status: %s\n", video.Status)
		return fmt.Errorf("Video creation failed. Status: %s\n", video.Status)
	}

	cr, err := t.data.OpenAI.Videos().DownloadContent(ctx, video.ID, openai.VideoDownloadContentParams{
		Variant: "video",
	})
	if err != nil {
		return err
	}

	content, err := io.ReadAll(cr.Body)
	if err != nil {
		return err
	}

	//videoz.GetFrame()
	//
	//url, err := t.data.Alioss.UploadBytes(ctx, helper.MD5(content)+".mp4", content)
	//if err != nil {
	//	return err
	//}

	url, err := t.data.TOS.Put(ctx, tos.PutRequest{
		Bucket:  "yoozyres",
		Content: content,
		Key:     helper.MD5(content) + ".mp4",
	})
	if err != nil {
		return err
	}

	_, err = t.data.Mongo.Asset.UpdateByIDXX(ctx, asset.XId, mgz.Set(bson.M{
		"coverUrl": tmpImages.Data[0].Url,
		"url":      url,
		"status":   "completed",
		"duration": conv.Float64(seconds),
		"extra": bson.M{
			"tmpImage": tmpImages.Data[0].Url,
		},
	}))
	if err != nil {
		return err
	}

	return nil
}
