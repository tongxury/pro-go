package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/wg"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjService) GenerateAISegment() {

	ctx := context.Background()

	list, err := t.data.Mongo.TaskSegment.List(ctx, bson.M{"status": "aiSegmentGenerating"})
	if err != nil {
		log.Errorw("List err", err)
		return
	}

	if len(list) == 0 {
		return
	}

	wg.WaitGroup(ctx, list, t.generateAISegment)
}

func (t ProjService) generateAISegment(ctx context.Context, taskSegment *projpb.TaskSegment) error {

	log.Debugw("start generateAISegment", taskSegment.XId)

	fastImages, errs := wg.WaitGroupResults(ctx, []int{1, 1, 1}, func(ctx context.Context, param int) (string, error) {

		image, err := t.doGenerateAISegmentV3(ctx, true, taskSegment)
		if err != nil {
			return "", err
		}

		return image, nil
	})

	if len(errs) != 0 {
		return errs[0]
	}

	lastImages, errs := wg.WaitGroupResults(ctx, []int{1, 1, 1}, func(ctx context.Context, param int) (string, error) {

		image, err := t.doGenerateAISegmentV3(ctx, false, taskSegment)
		if err != nil {
			return "", err
		}

		return image, nil
	})

	if len(errs) != 0 {
		return errs[0]
	}

	_, err := t.data.Mongo.TaskSegment.UpdateOneXXById(ctx, taskSegment.XId,
		bson.M{
			"frames.first": helper.Mapping(fastImages, func(x string) *projpb.Frame {
				return &projpb.Frame{
					Url: x,
				}
			}),
			"frames.last": helper.Mapping(lastImages, func(x string) *projpb.Frame {
				return &projpb.Frame{
					Url: x,
				}
			}),
			"status": "aiSegmentGenerated",
		})
	if err != nil {
		return err
	}

	return nil
}

func (t ProjService) doGenerateAISegmentV3(ctx context.Context, first bool, taskSegment *projpb.TaskSegment) (string, error) {

	log.Debugw("start doGenerateAISegmentV3", taskSegment.XId)

	//image := "https://yoozy.oss-cn-hangzhou.aliyuncs.com/0175810acead1db41b5d2894b126e2ba.jpg"
	//shangpin := "https://yoozy.oss-cn-hangzhou.aliyuncs.com/shangpin.png"

	shangpin0 := taskSegment.Task.Commodity.Images[0]
	shangpin1 := taskSegment.Task.Commodity.Images[1]
	shangpin2 := taskSegment.Task.Commodity.Images[2]
	origin := helper.Select(first, taskSegment.Frames.FirstOrigin, taskSegment.Frames.LastOrigin)

	params := model.GenerateImagesRequest{
		Model:  "doubao-seedream-4-0-250828",
		Prompt: "将图1的商品换为图2或图2图4中的商品, 并将图1中的文案去掉",
		Image: []string{
			origin,
			shangpin0,
			shangpin1,
			shangpin2,
		},
		Size:           volcengine.String("720x1280"),
		ResponseFormat: volcengine.String(model.GenerateImagesResponseFormatURL),
		//Watermark:      volcengine.Bool(true),
	}

	imagesResponse, err := t.data.Arkr.C().GenerateImages(ctx, params)
	if err != nil {
		log.Errorf("generate images error: %v\\n", err)
		return "", err
	}

	log.Debugw("doGenerateAISegmentV3", []string{
		origin,
		shangpin0,
		shangpin1,
		shangpin2,
		*imagesResponse.Data[0].Url,
	})

	return *imagesResponse.Data[0].Url, nil
}
