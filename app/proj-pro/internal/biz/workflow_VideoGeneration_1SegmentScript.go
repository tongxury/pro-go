package biz

import (
	"context"
	"fmt"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/data"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/third/gemini"

	"github.com/go-kratos/kratos/v2/log"
)

type VideoGeneration_SegmentScriptJob struct {
	data *data.Data
}

func (t VideoGeneration_SegmentScriptJob) Initialize(ctx context.Context, options Options) error {
	return nil
}

func (t VideoGeneration_SegmentScriptJob) GetName() string {
	return "segmentScriptJob"
}

func (t VideoGeneration_SegmentScriptJob) Execute(ctx context.Context, jobState *projpb.Job, wfState *projpb.Workflow) (status *ExecuteResult, err error) {

	dataBus := GetDataBus(wfState)

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "VideoGeneration_SegmentScriptJob.Execute",
		"workflowId ", wfState.XId,
		"jobState.Name ", jobState.Name,
		"jobState.Index ", jobState.Index,
		"jobState.Status ", jobState.Status,
	))

	logger.Debug("start")

	if dataBus.SegmentScript.GetScript() != "" {
		return &ExecuteResult{
			Status:      ExecuteStatusCompleted,
			SkipConfirm: true,
		}, nil
	}

	text, err := t.data.GenaiFactory.Get().GenerateContentV2(ctx, gemini.GenerateContentRequestV2{
		ImageUrls: dataBus.SegmentScript.Images,
		Prompt: fmt.Sprintf(`
请帮我根据我的灵感: %s, 生成一段视频脚本，用于生成一段视频， 不低于200字
`, dataBus.SegmentScript.GetInspiration()),
	})
	if err != nil {
		return nil, err
	}

	_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId, mgz.Op().
		Set(fmt.Sprintf("jobs.%d.dataBus.segmentScript.script", jobState.Index), text))

	if err != nil {
		logger.Errorw("update segment script fail", "err", err)
		return nil, err
	}

	return &ExecuteResult{
		Status: ExecuteStatusCompleted,
		//SkipConfirm: true,
	}, nil
}
