package biz

import (
	"context"
	"fmt"
	creditpb "store/api/credit"
	projpb "store/api/proj"
	"store/app/proj-pro/configs"
	"store/app/proj-pro/internal/data"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type WorkflowBiz struct {
	registry map[string]IWorkflow
	data     *data.Data
	//processing sync.Map
}

func NewWorkflowBiz(data *data.Data) *WorkflowBiz {

	videoReplication := NewVideoReplication(data)
	videoReplication2 := NewVideoReplication2(data)
	videoReplication3 := NewVideoReplication3(data)
	videoGeneration := NewVideoGeneration(data)

	t := &WorkflowBiz{
		registry: map[string]IWorkflow{
			videoReplication.GetName():  videoReplication,
			videoReplication2.GetName(): videoReplication2,
			videoReplication3.GetName(): videoReplication3,
			videoGeneration.GetName():   videoGeneration,
		},
		data: data,
	}

	go t.initRemixCache()

	return t
}

type CreatWorkFlowOptions struct {
	Auto bool
}

type Ops []CreatWorkFlowOptions

func (t Ops) GetAuto() bool {
	if len(t) == 0 {
		return false
	}

	return t[0].Auto
}

func (t *WorkflowBiz) createWorkflow(ctx context.Context, workflowName string, dataBus *projpb.DataBus, options ...CreatWorkFlowOptions) (*projpb.Workflow, error) {
	wfDef, ok := t.registry[workflowName]
	if !ok {
		return nil, fmt.Errorf("workflow %s not found", workflowName)
	}

	// Initialize Job States
	defs := wfDef.GetJobs()
	jobStates := make([]*projpb.Job, len(defs))
	for i, jobDef := range defs {
		if i == 0 {
			jobStates[i] = &projpb.Job{
				Name:   jobDef.GetName(),
				Status: JobStatusRunning,
				Index:  int64(i),
			}
		} else {
			jobStates[i] = &projpb.Job{
				Name:   jobDef.GetName(),
				Status: JobStatusWaiting,
				Index:  int64(i),
			}
		}
	}

	dataBus.Settings = &projpb.DataBus_Settings{
		AspectRatio: "9:16",
	}

	state := &projpb.Workflow{
		Name:          wfDef.GetName(),
		Current:       0,
		Status:        WorkflowStatusRunning,
		Jobs:          jobStates,
		DataBus:       dataBus,
		UserId:        dataBus.UserId,
		CreatedAt:     time.Now().Unix(),
		LastResumedAt: time.Now().Unix(),
		Auto:          Ops(options).GetAuto(),
	}

	insert, err := t.data.Mongo.Workflow.Insert(ctx, state)
	if err != nil {
		return nil, err
	}

	//err = t.data.Cache.Workflow.Set(ctx, &projpb.Workflow{XId: insert.XId}, 24*time.Hour)
	//if err != nil {
	//	return nil, err
	//}

	t.data.Redis.SAdd(ctx, "runningWorkflow:ids", insert.XId)

	return insert, nil
}

func (t *WorkflowBiz) Retry(ctx context.Context, workflowId string, jobIndex int64) error {

	workflow, err := t.data.Mongo.Workflow.GetById(ctx, workflowId)
	if err != nil {
		return err
	}

	for i := range workflow.Jobs {
		if int64(i) < jobIndex {
			continue
		}

		if int64(i) == jobIndex {

			workflow.Jobs[i].Status = JobStatusRunning
			workflow.Jobs[i].StartedAt = 0
			workflow.Jobs[i].CompletedAt = 0
			workflow.Jobs[i].DataBus = nil
		}

		if int64(i) > jobIndex {
			workflow.Jobs[i].Status = JobStatusWaiting
			workflow.Jobs[i].StartedAt = 0
			workflow.Jobs[i].CompletedAt = 0
			workflow.Jobs[i].DataBus = nil
		}
	}

	_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, workflowId,
		mgz.Op().Set("jobs", workflow.Jobs),
	)

	if err != nil {
		return err
	}

	return nil
}

func (t *WorkflowBiz) Back(ctx context.Context, workflowId string, jobIndex int64) error {

	workflow, err := t.data.Mongo.Workflow.GetById(ctx, workflowId)
	if err != nil {
		return err
	}

	for i := range workflow.Jobs {
		if int64(i) == jobIndex-1 {
			workflow.Jobs[i].Status = JobStatusConfirming
		}
		if int64(i) == jobIndex {
			workflow.Jobs[i] = &projpb.Job{
				Status: JobStatusWaiting,
				Name:   workflow.Jobs[i].Name,
				Index:  int64(i),
			}
		}
	}

	_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, workflowId,
		mgz.Op().
			Set("jobs", workflow.Jobs).
			Set("current", workflow.Current-1).
			Set("status", WorkflowStatusRunning),
	)

	if err != nil {
		return err
	}

	return nil
}

func (t *WorkflowBiz) Confirm(ctx context.Context, workflowId string, jobIndex int64, runImmediately bool) error {

	_, err := t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, workflowId,
		mgz.Op().SetListItem("jobs", int(jobIndex), "status", JobStatusCompleted))
	if err != nil {
		return err
	}

	workflow, err := t.data.Mongo.Workflow.GetById(ctx, workflowId)
	if err != nil {
		return err
	}

	jobDefs := t.registry[workflow.Name].GetJobs()

	nextJobDef := helper.SliceElement(jobDefs, int(jobIndex+1), false)
	if nextJobDef != nil {
		err = nextJobDef.Initialize(ctx, Options{
			JobState:       helper.SliceElement(workflow.Jobs, int(jobIndex+1), false),
			WorkflowState:  workflow,
			RunImmediately: runImmediately,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *WorkflowBiz) Cancel(ctx context.Context, workflowIds []string) error {

	_, err := t.data.Mongo.Workflow.UpdateByIDsIfExists(ctx, workflowIds,
		mgz.Op().
			Set("status", WorkflowStatusCanceled).
			Set("jobs.$[].status", JobStatusCanceled),
	)
	if err != nil {
		return err
	}

	return nil
}

func (t *WorkflowBiz) Pause(ctx context.Context, workflowIds []string) error {

	_, err := t.data.Mongo.Workflow.UpdateByIDsIfExists(ctx, workflowIds,
		mgz.Op().
			Set("status", WorkflowStatusPaused),
	)
	if err != nil {
		return err
	}

	return nil
}

func (t *WorkflowBiz) Resume(ctx context.Context, workflowIds []string) error {

	_, err := t.data.Mongo.Workflow.UpdateByIDsIfExists(ctx, workflowIds,
		mgz.Op().
			Set("status", WorkflowStatusRunning).
			Set("lastResumedAt", time.Now().Unix()),
	)
	if err != nil {
		return err
	}

	return nil
}

func (t *WorkflowBiz) Boost2(ctx context.Context, id string) error {

	//if id == "6969f2cfb48002731224605a" {
	//	fmt.Println(id)
	//}
	//
	//log.Debugw("boost task id", id)

	wfState, err := t.data.Mongo.Workflow.GetById(ctx, id)
	if err != nil {
		log.Errorw("GetById err ", err, "id", id)
		return err
	}

	return t.boost(ctx, wfState)
}

func (t *WorkflowBiz) boost(ctx context.Context, wfState *projpb.Workflow) error {

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "boost",
		"name", wfState.Name,
		"workflowId", wfState.XId,
	))

	//logger.Debug("boost workflow start")

	// LoadTasks Definition
	def, ok := t.registry[wfState.Name]
	if !ok {
		logger.Warnf("workflow %s not found", wfState.Name)
		return nil
	}

	jobDefs := def.GetJobs()

	// Safety check for index out of bounds
	if int(wfState.Current) >= len(jobDefs) {
		logger.Warnf("Workflow index out of bounds (Index: %d, Jobs: %d), marking completed", wfState.Current, len(jobDefs))
		//wfState.Status = WorkflowStatusCompleted
		//wfState.CompletedAt = time.Now().Unix()
		//_, _ = t.data.Mongo.Workflow.ReplaceByID(ctx, wfState.XId, wfState)

		t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
			mgz.Op().
				Set("status", WorkflowStatusCompleted).
				Set("completedAt", time.Now().Unix()),
		)

		return nil
	}

	if int(wfState.Current) >= len(wfState.Jobs) {

		newJobs := make([]*projpb.Job, len(jobDefs), len(jobDefs))
		for i, x := range jobDefs {

			if len(wfState.Jobs) > i {
				newJobs[i] = wfState.Jobs[i]
			} else {
				newJobs[i] = &projpb.Job{
					Name:      x.GetName(),
					Index:     int64(i),
					Status:    JobStatusWaiting,
					StartedAt: time.Now().Unix(),
				}
			}
		}

		_, err := t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
			mgz.Op().
				Set("jobs", newJobs).
				Set("status", WorkflowStatusRunning), // 没有意义
		)
		if err != nil {
			return err
		}
		return nil
	}

	currentJobDef := jobDefs[wfState.Current]
	currentJobState := wfState.Jobs[wfState.Current]
	//// Execution Logic
	//if currentJobState.Status != JobStatusConfirming {
	//	//logger.Debugw("checking job", "", "index", wfState.Current, "job", currentJobDef.GetName(), "status", currentJobState.Status)
	//}

	// 1. Handle Failed Job
	if currentJobState.Status == JobStatusFailed {
		//wfState.Status = WorkflowStatusFailed
		//wfState.CompletedAt = time.Now().Unix()
		//_, _ = t.data.Mongo.Workflow.ReplaceByID(ctx, wfState.XId, wfState)

		t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
			mgz.Op().
				Set("status", WorkflowStatusFailed).
				Set("completedAt", time.Now().Unix()),
		)

		return nil
	}

	if currentJobState.Status == JobStatusConfirming {
		return nil
	}

	// 2. Handle Completed Job (Transition to Next)
	if currentJobState.Status == JobStatusCompleted {
		// If this was the last job
		if int(wfState.Current) >= len(jobDefs)-1 {
			//wfState.Status = WorkflowStatusCompleted
			//wfState.CompletedAt = time.Now().Unix()

			err := def.OnComplete(ctx, wfState)
			if err != nil {
				logger.Errorw("OnComplete err", err)
				return nil
			}

			t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
				mgz.Op().
					Set("status", WorkflowStatusCompleted).
					Set("completedAt", time.Now().Unix()),
			)

		} else {
			// Move to next job
			logger.Debugw("OnComplete already completed")
			t.data.Mongo.Workflow.UpdateByIDIfExists(ctx,
				wfState.XId,
				mgz.Op().
					Set("current", currentJobState.Index+1),
			)

		}
		//_, _ = t.data.Mongo.Workflow.ReplaceByID(ctx, wfState.XId, wfState)
		// Process the next state immediately in the next Boost cycle (or could loop here)
		return nil
	}

	// 3. Handle Waiting Job (Prepare Input & Start)
	if currentJobState.Status == JobStatusWaiting {
		//currentJobState.StartedAt = time.Now().Unix()
		//currentJobState.Status = JobStatusRunning

		t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
			mgz.Op().
				Set(fmt.Sprintf("jobs.%d.startedAt", currentJobState.Index), time.Now().Unix()).
				Set(fmt.Sprintf("jobs.%d.status", currentJobState.Index), JobStatusRunning),
			//Set("status", WorkflowStatusFailed).
			//Set("completedAt", time.Now().Unix()),
		)
		return nil
		//// Update state before execution to prevent double execution if crash
		//_, _ = t.data.Mongo.Workflow.ReplaceByID(ctx, wfState.XId, wfState)
	}

	// 4. Execute Running Job
	if currentJobState.Status == JobStatusRunning {

		//if currentJobDef.IsComplete(ctx, currentJobState, wfState) {
		//
		//}

		// 加上回调函数 比如 onSave
		// 返回值要加上 其他状态 比如失败
		//t1 := time.Now()

		es, err := currentJobDef.Execute(ctx, currentJobState, wfState)

		//logger.Debugw("Execute ", "done", "duration", time.Since(t1))

		if err != nil {
			logger.Errorf("Job execution failed: %v", err)
			//currentJobState.Status = JobStatusFailed
			//wfState.Status = WorkflowStatusFailed
			//wfState.CompletedAt = time.Now().Unix()

			//t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
			//	mongoz.Op().
			//		Set(fmt.Sprintf("jobs.%d.status", currentJobState.Index), JobStatusFailed).
			//		Set("status", WorkflowStatusFailed).
			//		Set("completedAt", time.Now().Unix()),
			//)

			//_, _ = t.data.Mongo.Workflow.ReplaceByID(ctx, wfState.XId, wfState)
			return nil
		}

		if es == nil {
			return nil
		}

		if es.Status == ExecuteStatusCompleted {
			logger.Infof("Job completed: %s", currentJobDef.GetName())
			//currentJobState.Status = JobStatusCompleted
			//currentJobState.CompletedAt = time.Now().Unix()

			st := helper.Select(es.SkipConfirm || wfState.Auto, JobStatusCompleted, JobStatusConfirming)

			t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
				mgz.Op().
					Set(fmt.Sprintf("jobs.%d.status", currentJobState.Index), st).
					Set(fmt.Sprintf("jobs.%d.completedAt", currentJobState.Index), time.Now().Unix()),
				//Set("status", WorkflowStatusFailed).
				//Set("completedAt", time.Now().Unix()),
			)

			if es.Cost > 0 {

				_, err = t.data.GrpcClients.CreditClient.XCost(ctx, &creditpb.XCostRequest{
					UserId: wfState.UserId,
					Amount: configs.CreditCostAsset,
					Key:    wfState.XId,
				})
				if err != nil {
					logger.Errorw("Cost err", err, "userId", wfState.UserId)
				}

				logger.Debugw("XCost", "", "userId", wfState.UserId, "amount", configs.CreditCostAsset)

			}

			return nil
		}

		if es.Status == ExecuteStatusFailed {
			t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, wfState.XId,
				mgz.Op().
					Set(fmt.Sprintf("jobs.%d.status", currentJobState.Index), JobStatusFailed),
			)
			return nil
		}

		// If !ok, it means job is still running (async or long-running), do nothing.
	}
	return nil
}
