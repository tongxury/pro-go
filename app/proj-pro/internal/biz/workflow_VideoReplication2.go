package biz

import (
	"context"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/data"
	"store/pkg/clients/mgz"
)

type VideoReplication2 struct {
	Jobs []IJob
	data *data.Data
}

func (t *VideoReplication2) OnComplete(ctx context.Context, wfState *projpb.Workflow) error {

	dataBus := GetDataBus(wfState)
	vs := dataBus.VideoGenerations
	if len(vs) == 0 {
		return nil
	}

	_, err := t.data.Mongo.Asset.UpdateOne(ctx,
		mgz.Filter().EQ("workflow._id", wfState.XId).B(),
		mgz.Op().
			Set("coverUrl", vs[0].CoverUrl).
			Set("status", "completed").
			Set("url", vs[0].Url),
	)
	if err != nil {
		return err
	}

	return nil
}

func (t *WorkflowBiz) CreateVideoReplicationWorkflow2(ctx context.Context, userId string, segment *projpb.ResourceSegment, commodity *projpb.Commodity) (*projpb.Workflow, error) {

	return t.createWorkflow(ctx, "VideoReplication2", &projpb.DataBus{
		UserId:    userId,
		Segment:   segment,
		Commodity: commodity,
	})
}

func NewVideoReplication2(data *data.Data) *VideoReplication2 {

	// 不能随意组合  每个工作流不行重新开发  因为数据 不一样
	return &VideoReplication2{
		data: data,
		Jobs: []IJob{
			&VideoReplication2_SegmentScriptJob{data: data},
			&VideoReplication2_KeyFrameGenerationJob{data: data},
			&VideoReplication2_VideoGenerationJob{data: data},
		},
	}

}

func (t *VideoReplication2) GetName() string {
	return "VideoReplication2"
}

func (t *VideoReplication2) GetJobs() []IJob {
	return t.Jobs
}
