package service

import (
	"context"
	"store/app/proj-pro/internal/biz"
	"store/pkg/clients/cronexecutor"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (t *WorkFlowService) PauseTimeoutWorkflowIdsExecutor() cronexecutor.IExecutor {
	return cronexecutor.NewExecutor(func(ctx context.Context) ([]cronexecutor.ITask, error) {
		return []cronexecutor.ITask{
			cronexecutor.NewTask(time.Now().String(), func(ctx context.Context) error {

				log.Debug("clear timeout workflow ids executor start")

				//filter := mgz.Filter().
				//	EQ("status", biz.WorkflowStatusRunning).
				//	LT("lastResumedAt", time.Now().Add(-2*time.Second).Unix()).
				//	B()

				filter := bson.M{
					"status": biz.WorkflowStatusRunning,
					//"lastResumedAt": time.Now().Add(-2 * time.Second).Unix(),
					"$or": []bson.M{
						{"lastResumedAt": bson.M{"$exists": false}},
						{"lastResumedAt": bson.M{"$lt": time.Now().Add(-2 * time.Hour).Unix()}},
					},
				}

				timeoutWorkflows, err := t.data.Mongo.Workflow.List(ctx,
					filter,
				)
				if err != nil {
					log.Errorw("Failed to list running workflows", err)
					return err
				}

				log.Debugw("clear timeout workflow ids executor done", len(timeoutWorkflows))

				if len(timeoutWorkflows) == 0 {
					return nil
				}

				err = t.workflow.Pause(ctx, mgz.Ids(timeoutWorkflows))
				if err != nil {
					log.Errorw("Failed to cancel timeout workflow ids executor", err)
					return err
				}

				return nil
			}),
		}, nil
	})
}

func (t *WorkFlowService) SyncRunningWorkflowIdsExecutor() cronexecutor.IExecutor {
	return cronexecutor.NewExecutor(func(ctx context.Context) ([]cronexecutor.ITask, error) {
		return []cronexecutor.ITask{
			cronexecutor.NewTask(time.Now().String(), func(ctx context.Context) error {

				//log.Debug("sync running workflow ids executor start")

				runningWorkflows, err := t.data.Mongo.Workflow.List(ctx,
					mgz.Filter().EQ("status", biz.WorkflowStatusRunning).B(),
					mgz.Find().Paging(1, 50).B(),
				)
				if err != nil {
					log.Errorw("Failed to list running workflows", err)
					return err
				}

				log.Debugw("sync running workflow ids executor done", mgz.Ids(runningWorkflows))

				if len(runningWorkflows) == 0 {
					t.data.Redis.Del(ctx, "runningWorkflow:ids")
					return nil
				}

				t.data.Redis.Del(ctx, "runningWorkflow:ids")
				t.data.Redis.SAdd(ctx, "runningWorkflow:ids", mgz.Ids(runningWorkflows))
				return nil
			}),
		}, nil
	})
}

func (t *WorkFlowService) LoadTasks(ctx context.Context) ([]cronexecutor.ITask, error) {

	//runningWorkflows, err := t.data.Mongo.Workflow.List(ctx,
	//	mgz.Filter().EQ("status", biz.WorkflowStatusRunning).B(),
	//)
	//if err != nil {
	//	log.Errorw("Failed to list running workflows", err)
	//	return nil, err
	//}
	ids, err := t.data.Redis.SMembers(ctx, "runningWorkflow:ids").Result()
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return nil, nil
	}

	//log.Debugw("loading tasks ids", ids)

	return helper.Mapping(ids, func(x string) cronexecutor.ITask {
		return cronexecutor.NewTask(
			"workflow_"+x,
			func(ctx context.Context) error {
				return t.workflow.Boost2(ctx, x)
			})
	}), nil
}

//func (t *WorkFlowService) Run() {
//	t.workflow.Boost()
//}

//
//// 这个方法靠定时任务 每秒会执行一下 设计逻辑得时候要注意
//func (t *WorkFlowService) execute(ctx context.Context, job *projpb.Job, workflow *projpb.Workflow) error {
//
//	logger := log.NewHelper(log.With(log.DefaultLogger,
//		"func", "execute",
//		"job", job,
//	))
//
//	logger.Debugw("start", "")
//
//	switch job.Name {
//	case "generateVideoSegments":
//		if workflow.DataBus.GetSegment() == nil {
//			return errors.New("segment data empty")
//		}
//
//	case "replaceSegmentProductImage":
//		if workflow.DataBus.GetSegment() == nil {
//			return errors.New("segment data empty")
//		}
//
//		// 复制一份数据 以免覆盖
//
//		output := workflow.DataBus.GetReplaceSegmentProductImage().GetOutput()
//		//jobCtx := workflow.DataBus.GetReplaceSegmentProductImage().GetContext()
//
//		if output == nil {
//
//			//初始化
//			output = workflow.DataBus.GetSegment()
//
//			for i := range output.Segments {
//				output.Segments[i].StartFrame = ""
//				output.Segments[i].EndFrame = ""
//			}
//
//			t.data.Mongo.Workflow.UpdateByIDIfExists(ctx,
//				workflow.XId,
//				mongoz.Op().
//					Set("dataBus.replaceSegmentProductImage.output", output),
//			)
//			// 下次执行就有数据了
//			return nil
//		}
//
//		// 执行
//		for i := range output.Segments {
//			if output.Segments[i].StartFrame == "" {
//				imageBytes, err := t.data.GenaiFactory.Get().GenerateImage(ctx, gemini.GenerateImageRequest{
//					Images: []string{
//						workflow.DataBus.GetSegment().Segments[i].StartFrame,
//						workflow.DataBus.GetCommodity().GetMedias()[0].Url,
//					},
//					Prompt: "帮我将图1中的商品换成图2中的商品",
//				})
//				if err != nil {
//					return err
//				}
//
//				imageImageUrl, err := t.data.TOS.Put(ctx, tos.PutRequest{
//					Bucket:  "yoozyres",
//					Content: imageBytes,
//					Key:     helper.MD5(imageBytes) + ".jpg",
//				})
//				if err != nil {
//					return err
//				}
//
//				_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, workflow.XId, mongoz.Op().
//					Set(fmt.Sprintf("dataBus.replaceSegmentProductImage.output.segments.%d.startFrame", i), imageImageUrl))
//				if err != nil {
//					return err
//				}
//			}
//
//			if output.Segments[i].EndFrame == "" {
//				imageBytes, err := t.data.GenaiFactory.Get().GenerateImage(ctx, gemini.GenerateImageRequest{
//					Images: []string{
//						workflow.DataBus.GetSegment().Segments[i].EndFrame,
//						workflow.DataBus.GetCommodity().GetMedias()[0].Url,
//					},
//					Prompt: "帮我将图1中的商品换成图2中的商品",
//				})
//				if err != nil {
//					return err
//				}
//
//				imageImageUrl, err := t.data.TOS.Put(ctx, tos.PutRequest{
//					Bucket:  "yoozyres",
//					Content: imageBytes,
//					Key:     helper.MD5(imageBytes) + ".jpg",
//				})
//				if err != nil {
//					return err
//				}
//
//				_, err = t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, workflow.XId, mongoz.Op().
//					Set(fmt.Sprintf("dataBus.replaceSegmentProductImage.output.segments.%d.endFrame", i), imageImageUrl))
//				if err != nil {
//					return err
//				}
//			}
//
//		}
//
//		t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, workflow.XId, mongoz.Op().
//			Set(fmt.Sprintf("jobs.%d.status", job.Index), biz.JobStatusCompleted))
//
//		//newSegment := workflow.DataBus.GetNewSegment()
//
//		//helper.FindInStringMap(job.Extra.GetContext(), "")
//
//		//for _, x := range workflow.DataBus.GetTemplateDescription().GetSegments() {
//		//
//		//}
//
//	}
//
//	return nil
//}
//
//func (t *WorkFlowService) dispatch(ctx context.Context, x *projpb.Workflow) error {
//
//	jobCount := len(x.Jobs)
//	currentJob := x.Jobs[x.Current]
//
//	isLastJob := jobCount <= int(x.Current)+1
//
//	//
//	if currentJob.Status == biz.JobStatusWaiting {
//		_, err := t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, x.XId,
//			mongoz.Op().Sets(bson.M{
//				fmt.Sprintf("jobs.%d.status", x.Current):    biz.JobStatusRunning,
//				fmt.Sprintf("jobs.%d.startedAt", x.Current): time.Now().Unix()},
//			),
//		)
//		if err != nil {
//			log.Errorw("update workflow err", err)
//			return err
//		}
//	} else
//	//	running
//	if currentJob.Status == biz.JobStatusRunning {
//
//		err := t.execute(ctx, currentJob, x)
//		if err != nil {
//			log.Errorw("execute err", err)
//			return err
//		}
//	} else
//	// 已完成的job
//	if currentJob.Status == biz.JobStatusCompleted {
//		//最后1个job 继续推进
//		if isLastJob {
//			// 是最后1个job  结束整个workflow
//			_, err := t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, x.XId,
//				mongoz.Op().
//					Set(fmt.Sprintf("jobs.%d.completedAt", x.Current), time.Now().Unix()).
//					Set("status", biz.WorkflowStatusCompleted),
//			)
//			if err != nil {
//				log.Errorw("update workflow err", err)
//				return err
//			}
//		} else {
//			_, err := t.data.Mongo.Workflow.UpdateByIDIfExists(ctx, x.XId,
//				mongoz.Op().
//					Set("current", x.Current+1),
//			)
//			if err != nil {
//				log.Errorw("update workflow err", err)
//				return err
//			}
//		}
//	}
//
//	return nil
//}
