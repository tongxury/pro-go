package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/wg"
	"store/pkg/sdk/third/bytedance/volcengine"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (t *SessionService) Remix() {

	ctx := context.Background()

	list, err := t.data.Mongo.Session.List(ctx, bson.M{"status": "remixing"})
	if err != nil {
		log.Errorw("List err", err)
		return
	}

	if len(list) == 0 {
		return
	}

	wg.WaitGroup(ctx, list, t.remix)
}

func (t *SessionService) remix(ctx context.Context, session *projpb.Session) error {

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "remix",
		"item", session.XId,
	))

	logger.Debugw("start", "")

	taskId := helper.FindInStringMap(session.Extra.GetContext(), "taskId")

	if taskId != "" {

		result, err := t.data.Volcengine.QueryMixCutTaskResult(ctx, volcengine.QueryMixCutTaskResultParams{
			TaskKey: taskId,
		})
		if err != nil {
			logger.Errorw("Query MixCutTaskResult err", err)
			return err
		}

		if result.Data.Task.Status == 200 {

			//result.Data.Task.VideoList[0].DownloadUrl

			video, err := t.data.TOS.PutVideo(ctx, result.Data.Task.VideoList[0].DownloadUrl)
			if err != nil {

				logger.Errorw("PutVideo err", err)
				return err
			}

			coverUrl, err := t.data.TOS.PutImage(ctx, result.Data.Task.VideoList[0].CoverUrl)
			if err != nil {
				logger.Errorw("PutImage err", err)
				return err
			}

			t.data.Mongo.Session.UpdateByIDIfExists(ctx, session.XId,
				mgz.Op().
					Set("url", video).
					Set("coverUrl", coverUrl).
					Set("status", "completed"),
			)

			return nil
		}

		return nil
	}

	segments, err := t.data.Mongo.SessionSegment.List(ctx, mgz.Filter().EQ("sessionId", session.XId).B())
	if err != nil {
		return err
	}

	if len(segments) == 0 {
		return nil
	}

	var groups []volcengine.Group

	for _, x := range segments {

		assetId := x.Asset.GetXId()

		asset, err := t.data.Mongo.Asset.GetById(ctx, assetId)
		if err != nil {
			return err
		}

		material, err := t.data.Volcengine.CreateUrlMaterial(ctx, volcengine.CreateMaterialParams{
			MediaFirstCategory: "video",
			Title:              "title",
			MaterialUrl:        asset.Url,
			Wait:               true,
		})
		if err != nil {
			logger.Errorw("CreateUrlMaterial err", err, "url", asset.Url)
			return err
		}

		//x.Subtitle

		//parts := strings.Split(x.Subtitle, "，")

		//var texts []volcengine.TextConfig
		//for _, xx := range parts {
		//	for _, xxx := range stringz.SplitToSlice(xx, 10) {
		//		texts = append(texts, volcengine.TextConfig{
		//			Text: xxx,
		//
		//			Font: "7130161740481855501",
		//			//BgColor:   "#000000",
		//			FlowerId:     "7259353395813695544",
		//			FontSize:     "15",
		//			InnerPadding: 0.5,
		//
		//			//BgColor:   "#000000",
		//			Bold:      1,
		//			Alignment: "center",
		//			//X:         0.3,
		//			Y: 0.8,
		//		})
		//	}
		//}

		groups = append(groups, volcengine.Group{
			MuseList: []volcengine.Muse{
				{
					Type:      "video",
					MuseId:    material.MediaId,
					StartTime: 0,
					EndTime:   0,
				},
			},
			//TextConfig: []volcengine.TextConfig{
			//	{
			//		Text:      "视频混剪mix-cut1111&&&&&",
			//		Font:      "7028804401003380749",
			//		FontSize:  "14",
			//		BgColor:   "#FFC0CB",
			//		Alignment: "right",
			//		X:         0.5,
			//		Y:         0.8,
			//	},
			//	{
			//		Text:      "视频混剪测试1111",
			//		Font:      "7028804401003380749",
			//		FontSize:  "14",
			//		BgColor:   "#FFC0CB",
			//		Alignment: "right",
			//		X:         0.1,
			//		Y:         0.3,
			//	},
			//},
			//TtsConfig: []volcengine.TtsConfig{
			//	{
			//		//AddAudio:    1,
			//		//AddSubtitle: 1,
			//		//Subtitle:    texts,
			//		//TtsId: "12",
			//		ToneId: "33",
			//		Speed:  150,
			//	},
			//},
			//TimeStrategy: volcengine.TimeStrategy{
			//	Type: 1,
			//},
			Type:   2,
			Volume: 0,
		})
	}

	time.Sleep(10 * time.Second)
	cutTask, err := t.data.Volcengine.CreateMixCutTask(ctx, volcengine.CreateMixCutTaskAsyncParams{
		Name:             "混剪",
		Group:            groups,
		Height:           1920,
		Width:            1080,
		TransitionConfig: &volcengine.TransitionConfig{
			//Type: 1,
		},
		FilterConfig: &volcengine.FilterConfig{
			Type: 1,
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
	})

	if err != nil {
		logger.Errorw("create cut task err", err)
		return err
	}

	t.data.Mongo.Session.UpdateByIDIfExists(ctx, session.XId,
		mgz.Op().Set("extra.context.taskId", cutTask.TaskKey),
	)

	return nil
}
