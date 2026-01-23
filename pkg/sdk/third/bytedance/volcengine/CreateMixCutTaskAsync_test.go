package volcengine

import (
	"context"
	"fmt"
	"store/pkg/sdk/conv"
	"testing"
	"time"
)

//{"Code":0,"Message":"\u6210\u529f","ResponseMetadata":{"Action":"CreateMixCutTaskAsync","Region":"cn-north","RequestId":"20250926214913F346378D1DF3C46B5DE2","Service":"iccloud_muse","Version":"2025-03-26"},"Result":{"TaskKey":"9663242a-9adf-11f0-94f6-3436ac120034"}}

func TestClient_CreateMixCutTaskAsync(t *testing.T) {

	c := NewClient()

	ctx := context.Background()

	//async, err := c.CreateMixCutTaskAsync(ctx,
	//	CreateMixCutTaskAsyncParams{
	//		Name: "混剪API",
	//		Group: []Group{
	//			{
	//				Name: "视频组1",
	//				MuseList: []Muse{
	//					{
	//						Type:      "video",
	//						MuseId:    "7554445328678862894",
	//						StartTime: 0,
	//						EndTime:   3000,
	//					},
	//					{
	//						Type:      "video",
	//						MuseId:    "7554445328678862894",
	//						StartTime: 0,
	//						EndTime:   3000,
	//					},
	//				},
	//				TextConfig: []TextConfig{
	//					{
	//						Text:      "视频混剪mix-cut1111&&&&&",
	//						Font:      "7028804401003380749",
	//						FontSize:  "14",
	//						BgColor:   "#FFC0CB",
	//						Alignment: "right",
	//						X:         0.5,
	//						Y:         0.8,
	//					},
	//					{
	//						Text:      "视频混剪测试1111",
	//						Font:      "7028804401003380749",
	//						FontSize:  "14",
	//						BgColor:   "#FFC0CB",
	//						Alignment: "right",
	//						X:         0.1,
	//						Y:         0.3,
	//					},
	//				},
	//				TtsConfig: []TtsConfig{
	//					{
	//						AddAudio:    1,
	//						AddSubtitle: 1,
	//						Subtitle: []TextConfig{
	//							{
	//								Text:      "这是一段TTS语音",
	//								Font:      "7130161740481855501",
	//								FontSize:  "12",
	//								BgColor:   "#000000",
	//								FlowerId:  "7259353395813695544",
	//								Alignment: "left",
	//								X:         0.3,
	//								Y:         0.4,
	//							},
	//						},
	//						//TtsId: "12",
	//						ToneId: "33",
	//					},
	//				},
	//				Type:   2,
	//				Volume: 0,
	//			},
	//		},
	//		Height: 1920,
	//		Width:  1080,
	//		TransitionConfig: &TransitionConfig{
	//			Type: 1,
	//		},
	//		FilterConfig: &FilterConfig{
	//			Type: 1,
	//		},
	//		//Background: &Background{
	//		//	Type: 2,
	//		//	Rgba: "#ffffff00",
	//		//	Ids:  []string{"7483421790656561161"},
	//		//},
	//		//GlobalTextConfig: &GlobalTextConfig{
	//		//	Type: 1,
	//		//	Configs: [][]Config{
	//		//		{
	//		//			{
	//		//				Text:      "视频混剪mix-cut1111&&&&&",
	//		//				Font:      "7028804401003380749",
	//		//				FontSize:  "20",
	//		//				BgColor:   "#FFC0CB",
	//		//				Alignment: "right",
	//		//				X:         0.5,
	//		//				Y:         0.8,
	//		//			},
	//		//			{
	//		//				Text:      "视频混剪测试1111",
	//		//				Font:      "7028804401003380749",
	//		//				FontSize:  "20",
	//		//				BgColor:   "#FFC0CB",
	//		//				Alignment: "right",
	//		//				X:         0.1,
	//		//				Y:         0.3,
	//		//			},
	//		//		},
	//		//	},
	//		//},
	//	})
	//if err != nil {
	//	panic(err)
	//}
	//
	//time.Sleep(5 * time.Second)
	//preview, err := c.QueryMixCutTaskPreview(ctx, QueryMixCutTaskPreviewParams{
	//	TaskKey: async.TaskKey,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//
	//var groupIds []int
	//for _, x := range preview.Task {
	//	groupIds = append(groupIds, x.GroupId)
	//}
	//
	//submit, err := c.SubmitMixCutTaskAsync(ctx, SubmitMixCutTaskAsyncParams{
	//	TaskKey:  preview.TaskKey,
	//	GroupIds: groupIds,
	//})
	//if err != nil {
	//	panic(err)
	//}

	for {
		time.Sleep(1 * time.Second)
		result, err := c.QueryMixCutTaskResult(ctx, QueryMixCutTaskResultParams{
			TaskKey: "ba3a8064-9bb1-11f0-94f6-3436ac120034",
		})
		if err != nil {
			continue
		}

		if result.Data.Task.Status == 300 {
			fmt.Println(conv.S2J(result))
			break
		}

		if result.Data.Task.Status == 200 {
			fmt.Println(conv.S2J(result))
			fmt.Println(result.Data.Task.VideoList[0].DownloadUrl)
			break
		}

	}
}
