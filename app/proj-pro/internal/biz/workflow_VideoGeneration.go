package biz

import (
	"context"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/data"
	"store/pkg/clients/mgz"
)

type VideoGeneration struct {
	Jobs []IJob
	data *data.Data
}

func (t *VideoGeneration) OnComplete(ctx context.Context, wfState *projpb.Workflow) error {

	dataBus := GetDataBus(wfState)
	vs := dataBus.VideoGenerations
	if len(vs) == 0 {
		return nil
	}

	_, err := t.data.Mongo.Asset.UpdateOne(ctx,
		mgz.Filter().EQ("workflow._id", wfState.XId).B(),
		mgz.Op().
			Set("coverUrl", vs[0].CoverUrl).
			Set("url", vs[0].Url),
	)
	if err != nil {
		return err
	}

	return nil
}

func (t *WorkflowBiz) CreateVideoGenerationWorkflow(ctx context.Context, userId string, script *projpb.SegmentScript) (*projpb.Workflow, error) {

	return t.createWorkflow(ctx, "VideoGeneration", &projpb.DataBus{
		UserId:        userId,
		SegmentScript: script,
	})
}

func NewVideoGeneration(data *data.Data) *VideoGeneration {

	// 不能随意组合  每个工作流不行重新开发  因为数据 不一样
	return &VideoGeneration{
		data: data,
		Jobs: []IJob{
			&VideoGeneration_SegmentScriptJob{
				data: data,
			},
			&VideoGeneration_VideoGenerationJob{
				data: data,
			},
		},
	}

}

func (t *VideoGeneration) GetName() string {
	return "VideoGeneration"
}

func (t *VideoGeneration) GetJobs() []IJob {
	return t.Jobs
}
