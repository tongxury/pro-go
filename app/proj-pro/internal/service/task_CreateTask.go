package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/krathelper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjService) CreateTask(ctx context.Context, params *projpb.CreateTaskTaskRequest) (*projpb.Task, error) {

	userId := krathelper.RequireUserId(ctx)

	comm, err := t.data.Mongo.Commodity.FindByID(ctx, params.CommodityId)
	if err != nil {
		return nil, err
	}

	//change := comm.Chances[params.ChanceIndex]
	//change.Index = params.ChanceIndex

	//comm.Chances = nil

	newTask := &projpb.Task{
		Commodity:    comm,
		TargetChance: comm.Chances[0],
		Status:       "chanceSelecting",
		CreatedAt:    time.Now().Unix(),
		UserId:       userId,
		//Steps:     configs.InitialSteps,
	}

	//tags := append(comm.Tags, change.TargetAudience.GetTags()...)
	//for _, x := range change.SellingPoints {
	//	tags = append(tags, x.GetTags()...)
	//}

	task, _, err := t.data.Mongo.Task.InsertNX(ctx,
		newTask,
		bson.M{"userId": userId,
			"commodity._id":      params.CommodityId,
			"targetChance.index": params.ChanceIndex,
		},
	)
	if err != nil {
		return nil, err
	}

	newTask.XId = task

	return newTask, nil
}
