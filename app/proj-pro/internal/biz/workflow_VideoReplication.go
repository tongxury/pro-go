package biz

import (
	"context"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/data"
	"store/pkg/clients/mgz"
)

type VideoReplication struct {
	Jobs []IJob
	data *data.Data
}

func (t *VideoReplication) OnComplete(ctx context.Context, wfState *projpb.Workflow) error {
	dataBus := GetDataBus(wfState)
	_, err := t.data.Mongo.Asset.UpdateOne(ctx,
		mgz.Filter().EQ("workflow._id", wfState.XId).B(),
		mgz.Op().
			Set("coverUrl", dataBus.GetRemix().GetCoverUrl()).
			Set("status", "completed").
			Set("url", dataBus.GetRemix().GetUrl()),
	)
	if err != nil {
		return err
	}

	return nil
}

func (t *WorkflowBiz) CreateVideoReplicationWorkflow(ctx context.Context, userId string, name string, segment *projpb.ResourceSegment, commodity *projpb.Commodity, initialData *projpb.DataBus) (*projpb.Workflow, error) {

	if name == "" {
		name = "VideoReplication"
	}

	dataBus := initialData
	if dataBus == nil {
		dataBus = &projpb.DataBus{}
	}

	dataBus.Segment = segment
	dataBus.Commodity = commodity
	dataBus.UserId = userId

	return t.createWorkflow(ctx, name, dataBus)
}

func NewVideoReplication(data *data.Data) *VideoReplication {

	// 不能随意组合  每个工作流不行重新开发  因为数据 不一样
	return &VideoReplication{
		data: data,
		Jobs: []IJob{
			&VideoReplication_SegmentScriptJob{
				data: data,
			},
			&KeyFramesGenerationJob{
				data: data,
			},
			&VideoSegmentsGenerationJob{
				data: data,
			},
			&VideoSegmentsRemixJob{
				data: data,
			},
		},
	}

}

func (t *VideoReplication) GetName() string {
	return "VideoReplication"
}

func (t *VideoReplication) GetJobs() []IJob {
	return t.Jobs
}
