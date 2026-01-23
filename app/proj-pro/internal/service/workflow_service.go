package service

import (
	projpb "store/api/proj"
	"store/app/proj-pro/internal/biz"
	"store/app/proj-pro/internal/data"
)

type WorkFlowService struct {
	projpb.UnimplementedWorkflowServiceServer
	data     *data.Data
	workflow *biz.WorkflowBiz
}

func NewWorkflowService(data *data.Data, workflow *biz.WorkflowBiz) *WorkFlowService {
	return &WorkFlowService{
		data:     data,
		workflow: workflow,
	}
}
