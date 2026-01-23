package biz

import (
	"context"
	"errors"
	"fmt"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/data"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/third/gemini"

	"github.com/go-kratos/kratos/v2/log"
)

type SegmentScriptJob struct {
	data *data.Data
}

func (t SegmentScriptJob) Initialize(ctx context.Context, options Options) error {
	return nil
}

func (t SegmentScriptJob) GetName() string {
	return "segmentScriptJob"
}

func (t SegmentScriptJob) Execute(ctx context.Context, jobState *projpb.Job, wfState *projpb.Workflow) (status *ExecuteResult, err error) {

	dataBus := GetDataBus(wfState)

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "SegmentScriptJob.Execute",
		"workflowId ", wfState.XId,
		"jobState.Name ", jobState.Name,
		"jobState.Index ", jobState.Index,
		"jobState.Status ", jobState.Status,
	))

	segment := dataBus.Segment
	if segment == nil {
		return nil, errors.New("segment not found")
	}

	commodity := dataBus.Commodity
	if commodity == nil {
		return nil, errors.New("commodity not found")
	}

	if dataBus.SegmentScript.GetScript() != "" {
		return &ExecuteResult{
			Status: ExecuteStatusCompleted,
		}, nil
	}

	//settings, err := t.data.Mongo.Settings.FindOne(ctx, bson.M{})
	//if err != nil {
	//	logger.Errorw("Settings.FindOne err", err)
	//	return
	//}
	//
	//prompt := settings.VideoScript.GetContent()

	prompt, err := t.data.Mongo.Settings.GetPrompt(ctx, "sora2_video_prompt_generation")
	if err != nil {
		return nil, err
	}

	text, err := t.data.GenaiFactory.Get().GenerateContentV2(ctx, gemini.GenerateContentRequestV2{
		//ImageUrls: []string{dataBus.GetCommodity().GetMedias()[0].Url},
		Prompt: fmt.Sprintf(`
%s
===
要复刻的爆款脚本: %s
新商品信息: %s
`,
			prompt,
			dataBus.GetSegment().GetScript(),
			conv.S2J(dataBus.GetCommodity()),
		),
	})
	if err != nil {
		logger.Errorw("gen content error", err)
		return nil, err
	}

	//seg, err := videoz.GetSegmentByUrl(ctx, segment.Root.GetUrl(), segment.TimeStart, segment.TimeEnd)
	//if err != nil {
	//	return nil, err
	//}
	//
	//settings, err := t.data.Mongo.Settings.FindOne(ctx, bson.M{})
	//if err != nil {
	//	logger.Errorw("Settings.FindOne err", err)
	//	return
	//}
	//
	//prompt := settings.VideoScript.GetContent()
	//
	//video, err := t.data.GenaiFactory.Get().AnalyzeVideo(ctx, gemini.AnalyzeVideoRequest{
	//	VideoBytes: seg.Content,
	//	Prompt:     prompt,
	//})
	//if err != nil {
	//	logger.Errorw("t.data.GenaiFactory.AnalyzeVideo err", err)
	//	return nil, err
	//}
	//
	//if video == "" {
	//	return nil, nil
	//}

	_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId, mgz.Op().
		Set(fmt.Sprintf("jobs.%d.dataBus.segmentScript", jobState.Index), &projpb.SegmentScript{
			Script: text,
			//Segments: segments.Segments,
		}))
	if err != nil {
		logger.Errorw("update segment script fail", "err", err)
		return nil, err
	}

	//_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId, mgz.Op().
	//	Set("dataBus.segmentScript", &projpb.SegmentScript{
	//		Script: text,
	//	}))
	//if err != nil {
	//	return nil, err
	//}

	return &ExecuteResult{
		Status: ExecuteStatusCompleted,
	}, nil
}
