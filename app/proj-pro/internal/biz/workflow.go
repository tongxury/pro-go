package biz

import (
	"context"
	projpb "store/api/proj"
)

type IWorkflow interface {
	GetName() string
	GetJobs() []IJob
	OnComplete(ctx context.Context, wfState *projpb.Workflow) error
}

func GetDataBus(wfState *projpb.Workflow) *projpb.DataBus {
	res := &projpb.DataBus{}
	if wfState.DataBus != nil {
		res = wfState.DataBus
	}
	for _, job := range wfState.Jobs {
		if job.DataBus == nil {
			continue
		}
		if job.DataBus.Segment != nil {
			res.Segment = job.DataBus.Segment
		}
		if job.DataBus.Commodity != nil {
			res.Commodity = job.DataBus.Commodity
		}
		if len(job.DataBus.VideoGenerations) > 0 {
			res.VideoGenerations = job.DataBus.VideoGenerations
		}
		if len(job.DataBus.VideoFramesChanges) > 0 {
			res.VideoFramesChanges = job.DataBus.VideoFramesChanges
		}
		if job.DataBus.Remix != nil {
			res.Remix = job.DataBus.Remix
		}
		if job.DataBus.SegmentScript != nil {
			res.SegmentScript = job.DataBus.SegmentScript
		}
		if job.DataBus.KeyFrames != nil {
			res.KeyFrames = job.DataBus.KeyFrames
		}
	}
	return res
}
