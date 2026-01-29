package biz

import (
	"context"
	"fmt"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/data"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/stringz"
	"store/pkg/sdk/helper/videoz"
	"store/pkg/sdk/third/bytedance/volcengine"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type VideoReplication3_VideoSegmentsRemixJob struct {
	data *data.Data
}

func (t VideoReplication3_VideoSegmentsRemixJob) Initialize(ctx context.Context, options Options) error {

	wfState := options.WorkflowState
	jobState := options.JobState

	dataBus := GetDataBus(wfState)

	scriptSegments := dataBus.SegmentScript.GetSegments()
	if len(scriptSegments) == 0 {
		return fmt.Errorf("scriptSegments is empty")
	}

	videoSegments := dataBus.VideoGenerations
	if len(videoSegments) == 0 {
		return fmt.Errorf("videoSegments is empty")
	}

	segmentsRemix := dataBus.Remix
	// 初始化

	if segmentsRemix == nil {
		var remixSegments []*projpb.Remix_Segment

		//remixOptions, err := t.data.Cache.Remix.Get(ctx)
		//if err != nil {
		//	return err
		//}
		//if remixOptions == nil {
		//	return fmt.Errorf("remixOptions is empty")
		//}

		var scriptIndex int
		for _, x := range videoSegments {

			var ss *projpb.ResourceSegment
			if x.Category == "" {
				ss = scriptSegments[scriptIndex]
				scriptIndex += 1
			}

			remixSegments = append(remixSegments, &projpb.Remix_Segment{
				VideoGeneration: x,
				ScriptSegment:   ss,
				RemixOptions: &projpb.RemixOptions{
					//Flower: &projpb.RemixOptions_Flower{
					//	//DownloadUrl: "https://museapaas.aigc-cloud.com/api/storage/objects/media/7259353437698031675_origin.zip?infer_mime=ext&x-muse-token=ChYIARABGNzmh8sGINyJjcsGKgQ2ODUxEqYBCg4IroCBxrvFkNVoEKjFDxINCgt2b2xjX2VuZ2luZRoCGAIiCAoCaGwSAmNuKgYI2jQQhgMyKQoneyJiaXpfaWQiOiIxMDAwMDAwMDEiLCJkZXJpdmVkX3VpZCI6IiJ9QhQKEmljYXJjaC5pYW0uYXV0aHJwY0ouL2FwaS9zdG9yYWdlL29iamVjdHMvbWVkaWEvNzI1OTM1MzQzNzY5ODAzMTY3NRrYAnJ4cTZWMVA2Y3d5OFZyT1JEZHJNWURzTU5nQ21LSk10OXFDWlk4L3Y1cXl1V3FEeDNrM3hoNEtuVkczWmNlY0xoa1NvNzhGemI4TE9JUXVEaGlEekZQS3U3NkNtMkYweWR5aVdmSUZnNnUxcm5VaW5GcVpSTDYvZ0RyeTZyZnQwOEppREdFc0pqalB0QzJSMm8zK1dyNEg0NHhJQ2JmcW03b2UxWnJSdmZMaUkyWEo2QlFGNkZYOFZJeWRLNTRFMVR0Z1p2Q3BVaVNEa3pxREpVSHRiNEdMUEFLd1ArK2tUQ1pIa0VsbkhXZXhEUGM1aU9RYzZmM3g4b2RBYVFjZkhzM0Fzbjh0TEI1aFlaWDVELzVQVVQ1UjBrN29HbXE2bmlVdkQ3aUtZNGlCS2hMbFo1Yko5OEpxVFJIL0hUL09oeVFsakVVNFdTd1owV0VkSS9lYjQzQT09",
					//	MediaId: remixOptions.Flowers[0].MediaId,
					//	Cover:   remixOptions.Flowers[0].Cover,
					//	Name:    remixOptions.Flowers[0].Name,
					//},
					//Font: &projpb.RemixOptions_Font{
					//	//DownloadUrl: "https://museapaas.aigc-cloud.com/api/storage/objects/media/7264845232439558180_origin.ttf?infer_mime=ext&x-muse-token=ChYIARABGNPjh8sGINOGjcsGKgQ2ODUxEqYBCg4IroCBxrvFkNVoEKjFDxINCgt2b2xjX2VuZ2luZRoCGAIiCAoCbHESAmNuKgYI2jQQhgMyKQoneyJiaXpfaWQiOiIxMDAwMDAwMDEiLCJkZXJpdmVkX3VpZCI6IiJ9QhQKEmljYXJjaC5pYW0uYXV0aHJwY0ouL2FwaS9zdG9yYWdlL29iamVjdHMvbWVkaWEvNzI2NDg0NTIzMjQzOTU1ODE4MBrYAks0NzdiWE54TktSek1OcVNvemVHVU5CeUs4UVJNeEdXdDZ5WUpBVkoyaTFJVDZTdjBJTGtySmNMd1luY3BqSVhPVlBGdTBUby9HdEwzTmlqUHRMREZ0V2FZUHJCc05BRkVuOWlFQ3pOMnhRcU9xZnN4dThnVENTak1remVkZFIvQmYwY2QzNWloVmJjMnlWL0hwQlQxRlJ2SFB0d1hybTBjS3d5dnBPK2xJem1uTjV0ZXI5R0V4TXE5cmdkLzA2ZzRhU3dKMEx6WHkwdVljOFg2ZXdtenVseGZuYXJIN1RrOTlEeThhb2p6RjBJK1hEYjh5VVozcWJyTkcreDJtVVJIWnN4OE5KV1hmb0tZb0JjelRUbk5MUE1rakR5SmN4bGhVMlRRRTFDUFJwd0FRQzV1NFVhQnNuR2hLRXQweExRR1psR2djMktNcXlDOC9vZ2I3NFRZUT09",
					//	MediaId: remixOptions.Fonts[0].MediaId,
					//	Name:    remixOptions.Fonts[0].Name,
					//},
					//Tone: &projpb.RemixOptions_Tone{
					//	Id:          remixOptions.Tones[0].Id,
					//	Name:        remixOptions.Tones[0].Name,
					//	Description: remixOptions.Tones[0].Description,
					//	DownloadUrl: remixOptions.Tones[0].DownloadUrl,
					//},
				},
			})
		}

		segmentsRemix = &projpb.Remix{
			Segments: remixSegments,
			Status:   helper.Select(options.RunImmediately, ExecuteStatusRunning, ExecuteStatusWaiting),
		}

		_, err := t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
			mgz.Op().Set(fmt.Sprintf("jobs.%d.dataBus.remix", jobState.Index), segmentsRemix),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t VideoReplication3_VideoSegmentsRemixJob) GetName() string {
	return "videoSegmentsRemixJob"
}

//func (t VideoSegmentsRemixJob) Execute(ctx context.Context, jobState *projpb.Job, wfState *projpb.Workflow) (status *ExecuteResult, err error) {
//
//	logger := log.NewHelper(log.With(log.DefaultLogger,
//		"func", "VideoSegmentsRemixJob.Execute",
//		"workflowId ", wfState.XId,
//		"jobState.Name ", jobState.Name,
//		"jobState.Index ", jobState.Index,
//		"jobState.Status ", jobState.Status,
//	))
//
//	dataBus := GetDataBus(wfState)
//
//	scriptSegments := dataBus.SegmentScript.GetSegments()
//	if len(scriptSegments) == 0 {
//		return nil, fmt.Errorf("scriptSegments is empty")
//	}
//
//	videoSegments := dataBus.VideoGenerations
//	if len(videoSegments) == 0 {
//		return nil, fmt.Errorf("videoSegments is empty")
//	}
//
//	if len(videoSegments) < len(scriptSegments) {
//		return nil, fmt.Errorf("videoSegments less than scriptSegments ")
//	}
//
//	segmentsRemix := dataBus.Remix
//	// 初始化
//
//	if segmentsRemix == nil {
//		var remixSegments []*projpb.Remix_Segment
//
//		for _, x := range videoSegments {
//			remixSegments = append(remixSegments, &projpb.Remix_Segment{
//				VideoGeneration: x,
//				RemixOptions: &projpb.RemixOptions{
//					Flower: &projpb.RemixOptions_Flower{
//						//DownloadUrl: "https://museapaas.aigc-cloud.com/api/storage/objects/media/7259353437698031675_origin.zip?infer_mime=ext&x-muse-token=ChYIARABGNzmh8sGINyJjcsGKgQ2ODUxEqYBCg4IroCBxrvFkNVoEKjFDxINCgt2b2xjX2VuZ2luZRoCGAIiCAoCaGwSAmNuKgYI2jQQhgMyKQoneyJiaXpfaWQiOiIxMDAwMDAwMDEiLCJkZXJpdmVkX3VpZCI6IiJ9QhQKEmljYXJjaC5pYW0uYXV0aHJwY0ouL2FwaS9zdG9yYWdlL29iamVjdHMvbWVkaWEvNzI1OTM1MzQzNzY5ODAzMTY3NRrYAnJ4cTZWMVA2Y3d5OFZyT1JEZHJNWURzTU5nQ21LSk10OXFDWlk4L3Y1cXl1V3FEeDNrM3hoNEtuVkczWmNlY0xoa1NvNzhGemI4TE9JUXVEaGlEekZQS3U3NkNtMkYweWR5aVdmSUZnNnUxcm5VaW5GcVpSTDYvZ0RyeTZyZnQwOEppREdFc0pqalB0QzJSMm8zK1dyNEg0NHhJQ2JmcW03b2UxWnJSdmZMaUkyWEo2QlFGNkZYOFZJeWRLNTRFMVR0Z1p2Q3BVaVNEa3pxREpVSHRiNEdMUEFLd1ArK2tUQ1pIa0VsbkhXZXhEUGM1aU9RYzZmM3g4b2RBYVFjZkhzM0Fzbjh0TEI1aFlaWDVELzVQVVQ1UjBrN29HbXE2bmlVdkQ3aUtZNGlCS2hMbFo1Yko5OEpxVFJIL0hUL09oeVFsakVVNFdTd1owV0VkSS9lYjQzQT09",
//						MediaId: "7259353437698031675",
//						Cover:   "https://museapaas.aigc-cloud.com/api/storage/objects/media/7259436449911799819_origin.png?infer_mime=ext&x-muse-token=ChYIARABGKqFiMsGIKqojcsGKgQ2ODUxEqYBCg4IroCBxrvFkNVoEKjFDxINCgt2b2xjX2VuZ2luZRoCGAIiCAoCaGwSAmNuKgYI2jQQhgMyKQoneyJiaXpfaWQiOiIxMDAwMDAwMDEiLCJkZXJpdmVkX3VpZCI6IiJ9QhQKEmljYXJjaC5pYW0uYXV0aHJwY0ouL2FwaS9zdG9yYWdlL29iamVjdHMvbWVkaWEvNzI1OTQzNjQ0OTkxMTc5OTgxORrYAkw1WXdyd2ZUdDFsZXc4SG5HREl2Y1pGMDlDbi8wNUFRVVN3T0JTMnRlWjZpbXV6VTZZb3AzRUd6RnIxSEdoaUJBdmJBMUliN2hJQ0Y4NmZZb2hEcm5BbnRyMFhCdWt4a1dJTzhscXVNNDd4QWxHL3AzaytiVllaZ3lBYk1HMnE3eE5jcUhWQ2JKUUs0MmJpa25MU3dhaksrcUJNT0g1dE5GTXE4MEo4U0NSc0E2M1BzUkNpUWd2YTJxeUltbVloT1dlb00wZ2hNMEFJelZLKzJrRE9QeFM4U0IwSUVjckgzRjdlSDFsRGQ4YXpQaWlvYzBQSnBvVjdMRDVrbmd0cmI3VmZSMXRTZGVNRkhUUTJaUi9rd0FGMHRkWGE0Ky9aMWpUWHVMcVlaanFZaUJQYXRIZURtaERCWElDdmxyVG9TMU9wa2w2dkRob05CNW1DdXV1K2Zvdz09",
//						//"DowloadUrl" : "https://museapaas.aigc-cloud.com/api/storage/objects/media/7259353558347300904_origin.zip?infer_mime=ext&x-muse-token=ChYIARABGKqFiMsGIKqojcsGKgQ2ODUxEqYBCg4IroCBxrvFkNVoEKjFDxINCgt2b2xjX2VuZ2luZRoCGAIiCAoCaGwSAmNuKgYI2jQQhgMyKQoneyJiaXpfaWQiOiIxMDAwMDAwMDEiLCJkZXJpdmVkX3VpZCI6IiJ9QhQKEmljYXJjaC5pYW0uYXV0aHJwY0ouL2FwaS9zdG9yYWdlL29iamVjdHMvbWVkaWEvNzI1OTM1MzU1ODM0NzMwMDkwNBrYAnR6UUZSc0hBK0pVT0tNdzRBTkR1TmJRRUtiSnJ0emtqdFMvcXJzeWtaUExzQmVQR3JsNFBIN2Qxd3pFNStlYmZyL1hqL2dzTHRrVnN3TTRYYnlJT2NUam05Nm1ibTB5bWRVRnFUQ2hKamluNHY5R0loWGFWbnpEZFZESUZqMUZ2WVdPbmY2QjRSMURoZEFQdTZCMFpzRnlCVUNlYlIzTFpsRE4vK3g4TFhYY29HWERXUHd6ekdMWXg2a0NuNTROQUpvR3luZWg2aGpHSmZNS0JkbjNuN3JPM2hlOXNOdEhOTnhxNlBVYmVOdGJ2ZzdIZzUya0R6dWo0S0x0NCtXVVhRN0VsZTVmdTZZME9jRFR6SWo0QkF4aENBcThoTHppUlJJd3ZHOTVMRlMxaUxieEQwMW5STVZSd04xR1d4ZFRoSGVQbDRNZTF5bSt6OFlHemEySzFrZz09",
//						//"MediaId" : "7259353558347300904",
//						Name: "元旦-电影2",
//					},
//					Font: &projpb.RemixOptions_Font{
//						//DownloadUrl: "https://museapaas.aigc-cloud.com/api/storage/objects/media/7264845232439558180_origin.ttf?infer_mime=ext&x-muse-token=ChYIARABGNPjh8sGINOGjcsGKgQ2ODUxEqYBCg4IroCBxrvFkNVoEKjFDxINCgt2b2xjX2VuZ2luZRoCGAIiCAoCbHESAmNuKgYI2jQQhgMyKQoneyJiaXpfaWQiOiIxMDAwMDAwMDEiLCJkZXJpdmVkX3VpZCI6IiJ9QhQKEmljYXJjaC5pYW0uYXV0aHJwY0ouL2FwaS9zdG9yYWdlL29iamVjdHMvbWVkaWEvNzI2NDg0NTIzMjQzOTU1ODE4MBrYAks0NzdiWE54TktSek1OcVNvemVHVU5CeUs4UVJNeEdXdDZ5WUpBVkoyaTFJVDZTdjBJTGtySmNMd1luY3BqSVhPVlBGdTBUby9HdEwzTmlqUHRMREZ0V2FZUHJCc05BRkVuOWlFQ3pOMnhRcU9xZnN4dThnVENTak1remVkZFIvQmYwY2QzNWloVmJjMnlWL0hwQlQxRlJ2SFB0d1hybTBjS3d5dnBPK2xJem1uTjV0ZXI5R0V4TXE5cmdkLzA2ZzRhU3dKMEx6WHkwdVljOFg2ZXdtenVseGZuYXJIN1RrOTlEeThhb2p6RjBJK1hEYjh5VVozcWJyTkcreDJtVVJIWnN4OE5KV1hmb0tZb0JjelRUbk5MUE1rakR5SmN4bGhVMlRRRTFDUFJwd0FRQzV1NFVhQnNuR2hLRXQweExRR1psR2djMktNcXlDOC9vZ2I3NFRZUT09",
//						MediaId: "7264845232439558180",
//						Name:    "抖音美好体",
//					},
//					Tone: &projpb.RemixOptions_Tone{
//						Id:          "2",
//						Name:        "成熟女声",
//						Description: "女声",
//						DownloadUrl: "https://lf-iccloud-muse.volcmusecdn.com/obj/labcv-tob/muse/new_tts_BV009_DPE_streaming_v2.mp3",
//					},
//				},
//			})
//		}
//
//		segmentsRemix = &projpb.Remix{
//			Segments: remixSegments,
//			Status:   ExecuteStatusWaiting,
//		}
//
//		t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
//			mgz.Op().Set(fmt.Sprintf("jobs.%d.dataBus.remix", jobState.Index), segmentsRemix),
//		)
//	}
//
//	if segmentsRemix.Status == ExecuteStatusRunning {
//
//		var concats []videoz.ConcatVideoSegment
//		for _, x := range segmentsRemix.Segments {
//
//			videoUrl := x.VideoGeneration.GetUrl()
//			speed := x.RemixOptions.GetSpeed()
//			timeStart := x.RemixOptions.GetTimeStart()
//			timeEnd := x.RemixOptions.GetTimeEnd()
//
//			source, err := videoz.GetBytes(videoUrl)
//			if err != nil {
//				return nil, err
//			}
//
//			concats = append(concats, videoz.ConcatVideoSegment{
//				Source:    source,
//				Speed:     speed,
//				TimeStart: timeStart,
//				TimeEnd:   timeEnd,
//			})
//		}
//
//		newVideoBytes, err := videoz.ConcatVideos(concats)
//		if err != nil {
//			logger.Errorw("ConcatVideos err", err)
//			return nil, err
//		}
//
//		newVideoUrl, err := t.data.TOS.PutVideoBytes(ctx, newVideoBytes)
//		if err != nil {
//			return nil, err
//		}
//
//		t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
//			mgz.Op().
//				Set(fmt.Sprintf("jobs.%d.dataBus.remix.url", jobState.Index), newVideoUrl).
//				Set(fmt.Sprintf("jobs.%d.dataBus.remix.status", jobState.Index), ExecuteStatusCompleted),
//		)
//
//		return &ExecuteResult{
//			Status: ExecuteStatusCompleted,
//		}, nil
//	}
//
//	return nil, nil
//}

func (t VideoReplication3_VideoSegmentsRemixJob) Execute(ctx context.Context, jobState *projpb.Job, wfState *projpb.Workflow) (status *ExecuteResult, err error) {

	dataBus := GetDataBus(wfState)

	segmentsRemix := dataBus.Remix
	// 初始化

	if segmentsRemix == nil {

		err = t.Initialize(ctx, Options{
			JobState:       jobState,
			WorkflowState:  wfState,
			RunImmediately: wfState.Auto,
		})
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	// task
	if segmentsRemix.GetTaskId() != "" {

		result, err := t.data.Volcengine.QueryMixCutTaskResult(ctx, volcengine.QueryMixCutTaskResultParams{
			TaskKey: segmentsRemix.GetTaskId(),
		})
		if err != nil {
			return nil, err
		}

		if result.Data.Task.Status == 200 {

			//result.Data.Task.VideoList[0].DownloadUrl

			video, err := t.data.TOS.PutVideo(ctx, result.Data.Task.VideoList[0].DownloadUrl)
			if err != nil {
				return nil, err
			}

			coverUrl, err := t.data.TOS.PutImage(ctx, result.Data.Task.VideoList[0].CoverUrl)
			if err != nil {
				return nil, err
			}

			if video != "" {
				t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
					mgz.Op().
						Set(fmt.Sprintf("jobs.%d.dataBus.remix.url", jobState.Index), video).
						Set(fmt.Sprintf("jobs.%d.dataBus.remix.status", jobState.Index), ExecuteStatusCompleted).
						Set(fmt.Sprintf("jobs.%d.dataBus.remix.coverUrl", jobState.Index), coverUrl),
				)

				return &ExecuteResult{
					Status: ExecuteStatusCompleted,
				}, nil
			}
		}

		return nil, nil
	}

	if segmentsRemix.Status == ExecuteStatusRunning {

		//
		var groups []volcengine.Group

		for _, x := range segmentsRemix.Segments {

			videoUrl := x.VideoGeneration.GetUrl()
			speed := x.RemixOptions.GetSpeed()
			timeStart := x.RemixOptions.GetTimeStart()
			timeEnd := x.RemixOptions.GetTimeEnd()

			if (speed > 0 && speed != 1.0) || timeStart > 0 || timeEnd > 0 {
				source, err := videoz.GetBytes(videoUrl)
				if err != nil {
					return nil, err
				}

				if speed == 0 {
					speed = 1.0
				}

				edited, err := videoz.EditVideo(source, videoz.EditParams{
					From:  timeStart,
					To:    timeEnd,
					Speed: speed,
				})
				if err != nil {
					return nil, err
				}

				newVideoUrl, err := t.data.TOS.PutVideoBytes(ctx, edited)
				if err != nil {
					return nil, err
				}
				videoUrl = newVideoUrl
			}

			material, err := t.data.Volcengine.CreateUrlMaterial(ctx, volcengine.CreateMaterialParams{
				MediaFirstCategory: "video",
				Title:              "title",
				MaterialUrl:        videoUrl,
				Wait:               true,
			})
			if err != nil {
				log.Errorw("CreateUrlMaterial err", err, "url", x.VideoGeneration.GetUrl())
				return nil, err
			}

			var tts []volcengine.TtsConfig
			var timeStrategy int

			font := x.RemixOptions.GetFont()
			flower := x.RemixOptions.GetFlower()
			tone := x.RemixOptions.GetTone()

			if x.ScriptSegment != nil && (font != nil || flower != nil || tone != nil) {
				//x.Subtitle
				parts := strings.Split(x.ScriptSegment.Subtitle, "，")
				var texts []volcengine.TextConfig

				for _, xx := range parts {
					for _, xxx := range stringz.SplitToSlice(xx, 20) {
						texts = append(texts, volcengine.TextConfig{
							Text: xxx,

							Font: font.GetMediaId(),
							//BgColor:   "#000000",
							FlowerId:     flower.GetMediaId(),
							FontSize:     "15",
							InnerPadding: 0.5,

							//BgColor:   "#000000",
							Bold:      1,
							Alignment: "center",
							//X:         0.3,
							Y: 0.8,
						})
					}
				}

				tts = []volcengine.TtsConfig{
					{
						AddAudio:    1,
						AddSubtitle: 1,
						Subtitle:    texts,
						//TtsId: "12",
						ToneId: tone.GetId(),
						Speed:  150,
					},
				}

				timeStrategy = 1
			}

			groups = append(groups, volcengine.Group{
				MuseList: []volcengine.Muse{
					{
						Type:      "video",
						MuseId:    material.MediaId,
						StartTime: 0,
						EndTime:   0,
					},
				},

				TtsConfig: tts,
				TimeStrategy: volcengine.TimeStrategy{
					Type: timeStrategy,
				},
				Type:   2,
				Volume: 0,
			})
		}

		//time.Sleep(10 * time.Second)
		params := volcengine.CreateMixCutTaskAsyncParams{
			Name:   "混剪",
			Group:  groups,
			Height: 1920,
			Width:  1080,
			TransitionConfig: &volcengine.TransitionConfig{
				//Type: 1,
			},
			FilterConfig: &volcengine.FilterConfig{
				Type: 0,
			},
			//Background: &Background{
			//	Type: 2,
			//	Rgba: "#ffffff00",
			//	Ids:  []string{"7483421790656561161"},
			//},
			//GlobalTextConfig: &GlobalTextConfig{
			//	Type: 1,
			//	Configs: [][]Config{
			//		{
			//			{
			//				Text:      "视频混剪mix-cut1111&&&&&",
			//				Font:      "7028804401003380749",
			//				FontSize:  "20",
			//				BgColor:   "#FFC0CB",
			//				Alignment: "right",
			//				X:         0.5,
			//				Y:         0.8,
			//			},
			//			{
			//				Text:      "视频混剪测试1111",
			//				Font:      "7028804401003380749",
			//				FontSize:  "20",
			//				BgColor:   "#FFC0CB",
			//				Alignment: "right",
			//				X:         0.1,
			//				Y:         0.3,
			//			},
			//		},
			//	},
			//},
		}
		cutTask, err := t.data.Volcengine.CreateMixCutTask(ctx, params)
		if err != nil {
			log.Errorw("create cut task err", err, "params", conv.S2J(params), "cutTask", cutTask)
			return nil, err
		}

		t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
			mgz.Op().
				Set(fmt.Sprintf("jobs.%d.dataBus.remix.taskId", jobState.Index), cutTask.TaskKey).
				Set(fmt.Sprintf("jobs.%d.dataBus.remix.status", jobState.Index), ExecuteStatusRunning),
		)

	}

	return nil, nil
}
