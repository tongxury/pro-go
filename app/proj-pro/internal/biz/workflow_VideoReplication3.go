package biz

import (
	"context"
	projpb "store/api/proj"
	"store/app/proj-pro/internal/data"
)

type VideoReplication3 struct {
	Jobs []IJob
	data *data.Data
}

func (t *VideoReplication3) OnComplete(ctx context.Context, wfState *projpb.Workflow) error {
	return nil
}

func (t *WorkflowBiz) CreateVideoReplication3Workflow(ctx context.Context, userId string, initialData *projpb.DataBus) (*projpb.Workflow, error) {

	dataBus := initialData
	if dataBus == nil {
		dataBus = &projpb.DataBus{}
	}

	dataBus.UserId = userId
	return t.createWorkflow(ctx, "VideoReplication3", dataBus, CreatWorkFlowOptions{Auto: true})
}

func NewVideoReplication3(data *data.Data) *VideoReplication3 {

	// 不能随意组合  每个工作流不行重新开发  因为数据 不一样
	return &VideoReplication3{
		data: data,
		Jobs: []IJob{
			&VideoReplication3_CommodityAnalysisJob{
				data: data,
			},
			&VideoReplication3_SegmentScriptJob{
				data: data,
			},
			&VideoReplication3_KeyFramesGenerationJob{
				data: data,
			},
			//&VideoReplication3_VideoSegmentsGenerationJob{
			//	data: data,
			//},
			//&VideoReplication3_VideoSegmentsRemixJob{
			//	data: data,
			//},
		},
	}

}

func (t *VideoReplication3) GetName() string {
	return "VideoReplication3"
}

func (t *VideoReplication3) GetJobs() []IJob {
	return t.Jobs
}
