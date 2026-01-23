package service

import (
	"context"
	"fmt"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper/imagez"
	"store/pkg/sdk/helper/wg"
	"store/pkg/sdk/third/gemini"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (t *SessionService) GenerateNewKeyFrames() {

	ctx := context.Background()

	list, err := t.data.Mongo.SessionSegment.List(ctx, bson.M{"status": "subtitleGenerated"})
	if err != nil {
		log.Errorw("List err", err)
		return
	}

	if len(list) == 0 {
		return
	}

	wg.WaitGroup(ctx, list, t.generateNewKeyFrames)
}

func (t *SessionService) generateNewKeyFrames(ctx context.Context, sessionSegment *projpb.SessionSegment) error {

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "generateNewKeyFrames",
		"item", sessionSegment.XId,
	))

	logger.Debugw("start generate newKeyFrames")
	res, err2 := t.generateKeyFramesByGenai(ctx, sessionSegment)
	if err2 != nil {
		return err2
	}

	_, err := t.data.Mongo.SessionSegment.UpdateByIDIfExists(ctx,
		sessionSegment.XId,
		mgz.Op().Sets(
			bson.M{
				"keyFrames": res,
				"status":    "keyFramesGenerated",
			}),
	)
	if err != nil {
		return err
	}
	return nil
}

func (t *SessionService) generateKeyFramesByGenai(ctx context.Context, sessionSegment *projpb.SessionSegment) (*projpb.KeyFrames, error) {
	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "generateKeyFramesByGenai",
		"item", sessionSegment.XId,
	))

	commodity := sessionSegment.Session.Commodity

	p := fmt.Sprintf(`
基础设置：
	3x3 无缝网格布局，9:16 纵横比，8K 分辨率。
风格：
	抖音（Douyin）原生爆款视频风格，由 iPhone 15 Pro 拍摄，原片直出质感（Raw footage aesthetic），自然的皮肤纹理，真实光影，手持摄影感。
**产品还原度**：
	[产品名称] 必须与提供给你的参考图 1:1 精准还原，包括文字、Logo 和包装细节，绝无变形。
分镜细节：
	%s
技术参数：
	无网格线，无边框，无缝拼接，超写实，高动态范围（HDR），真实的电商直播氛围。 --ar 9:16 --v 6.0
=== 
附加信息:
- 模仿的电商视频的关键信息为: %s,
	`,
		sessionSegment.Segment.Script,
		conv.S2J(&projpb.ResourceSegment{
			Formula:       sessionSegment.Segment.Formula,
			TypedTags:     sessionSegment.Segment.TypedTags,
			Style:         sessionSegment.Segment.Style,
			SceneStyle:    sessionSegment.Segment.SceneStyle,
			ContentStyle:  sessionSegment.Segment.ContentStyle,
			ShootingStyle: sessionSegment.Segment.ShootingStyle,
		}))

	blob, err := t.data.GenaiFactory.Get().GenerateImage(ctx, gemini.GenerateImageRequest{
		Images: []string{commodity.Medias[0].GetUrl()},
		//Videos: [][]byte{seg.Content},
		Prompt: p,
		//Count: 8,
	})
	if err != nil {
		logger.Errorw("GenerateImageV2 err", err)
		return nil, err
	}
	//
	//tmpUrl, err := t.data.TOS.PutImageBytes(ctx, blob)
	//if err != nil {
	//	return nil, err
	//}
	//
	//fmt.Println("tmpUrl", tmpUrl)

	images, err := imagez.Split3x3(blob)
	if err != nil {
		logger.Errorw("Split3x3 err", err)
		return nil, err
	}

	var frames []*projpb.KeyFrames_Frame

	for _, x := range images {

		tmpUrl, err := t.data.TOS.PutImageBytes(ctx, x)
		if err != nil {
			return nil, err
		}

		frames = append(frames, &projpb.KeyFrames_Frame{
			Url: tmpUrl,
		})
	}

	return &projpb.KeyFrames{
		Frames: frames,
	}, nil
}
