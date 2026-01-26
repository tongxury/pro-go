package service

import (
	projpb "store/api/proj"
	"store/app/proj-pro/internal/biz"
	"store/app/proj-pro/internal/data"
)

type ProjService struct {
	projpb.UnimplementedProjProServiceServer
	data     *data.Data
	workflow *biz.WorkflowBiz
}

func NewProjService(data *data.Data, workflow *biz.WorkflowBiz) *ProjService {
	s := &ProjService{
		data:     data,
		workflow: workflow,
	}

	s.Start()

	return s
}
