package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper/stringz"
	"store/pkg/sdk/helper/wg"
	"store/pkg/sdk/third/bytedance/volcengine"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjService) GenerateVideo() {

	ctx := context.Background()

	list, err := t.data.Mongo.Task.List(ctx, bson.M{"status": "videoGenerating"})
	if err != nil {
		log.Errorw("List err", err)
		return
	}

	if len(list) == 0 {
		return
	}

	wg.WaitGroup(ctx, list, t.generateVideo)
}

func (t ProjService) generateVideo(ctx context.Context, task *projpb.Task) error {

	log.Debugw("generateVideo", task.XId)

	segments, err := t.data.Mongo.TaskSegment.List(ctx, bson.M{"taskId": task.XId})
	if err != nil {
		return err
	}

	var groups []volcengine.Group

	for _, x := range segments {

		asset, err := t.data.Mongo.Asset.GetById(ctx, x.AssetId)
		if err != nil {
			return err
		}

		//targets := helper.Filter(x.GeneratedTasks, func(param *projpb.GeneratedTask) bool {
		//	return param.Selected
		//})
		//
		//y := x.GeneratedTasks[0]
		//if len(targets) != 0 {
		//	y = targets[0]
		//}

		log.Debugw("generateVideo", task.XId)

		material, err := t.data.Volcengine.CreateUrlMaterial(ctx, volcengine.CreateMaterialParams{
			MediaFirstCategory: "video",
			Title:              x.XId,
			MaterialUrl:        asset.Url,
			Wait:               true,
		})
		if err != nil {
			log.Errorw("CreateUrlMaterial err", err, "asset", asset.XId)
			return err
		}

		log.Debugw("generateVideo", task.XId, "material", material.MediaId)

		//x.Subtitle

		parts := strings.Split(x.Subtitle, "，")

		var texts []volcengine.TextConfig
		for _, xx := range parts {
			for _, xxx := range stringz.SplitToSlice(xx, 10) {
				texts = append(texts, volcengine.TextConfig{
					Text: xxx,

					Font: "7130161740481855501",
					//BgColor:   "#000000",
					FlowerId:     "7259353395813695544",
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
			TtsConfig: []volcengine.TtsConfig{
				{
					AddAudio:    1,
					AddSubtitle: 1,
					Subtitle:    texts,
					//TtsId: "12",
					ToneId: "33",
					Speed:  150,
				},
			},
			TimeStrategy: volcengine.TimeStrategy{
				Type: 1,
			},
			Type:   2,
			Volume: 0,
		})
	}

	time.Sleep(10 * time.Second)
	cutTask, err := t.data.Volcengine.CreateMixCutTask(ctx, volcengine.CreateMixCutTaskAsyncParams{
		Name:   "混剪",
		Group:  groups,
		Height: 1920,
		Width:  1080,
		TransitionConfig: &volcengine.TransitionConfig{
			Type: 1,
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
		log.Errorw("create cut task err", task.XId, "err", err)
		return err
	}

	_, err = t.data.Mongo.Task.UpdateByIDIfExists(ctx,
		task.XId,
		mgz.Op().
			Set("generatedResult.taskId", cutTask.TaskKey).
			Set("status", "videoGeneratingWaiting"),
	)

	if err != nil {
		return err
	}

	return nil
}
