package service

import (
	"context"
	"fmt"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjService) UpdateTask(ctx context.Context, params *projpb.UpdateTaskRequest) (*projpb.Task, error) {

	switch params.Action {
	case "updateMode":

		list, err := t.data.Mongo.TaskSegment.List(ctx, bson.M{"taskId": params.Id})
		if err != nil {
			return nil, err
		}

		for _, x := range list {
			if x.AssetId != "" {
				fmt.Println(x.AssetId)
				return nil, errors.Forbidden("", "")
			}
		}

		_, err = t.data.Mongo.Task.UpdateByIDIfExists(ctx, params.Id,
			mgz.Op().
				Set("mode", params.Params["mode"]),
		)
		if err != nil {
			log.Error(err)
			return nil, err
		}

	case "updateTargetChance":

		task, err := t.data.Mongo.Task.GetById(ctx, params.Id)
		if err != nil {
			return nil, err
		}

		chance := task.Commodity.Chances[int(params.ChanceIndex)]

		_, err = t.data.Mongo.Task.UpdateByIDIfExists(ctx,
			params.Id,
			mgz.Op().
				Set("status", "templateSelecting").
				Set("targetChance", chance),
		)

		if err != nil {
			log.Errorw("UpdateOneXXById err", err)
			return nil, err
		}
	case "updateTemplate":

		task, err := t.data.Mongo.Task.GetById(ctx, params.Id)
		if err != nil {
			return nil, err
		}

		template, err := t.data.GrpcClients.ProjAdminClient.GetTemplate(ctx, &projpb.GetGetTemplateRequest{
			Id: params.TemplateId,
		})
		if err != nil {
			return nil, err
		}

		if template == nil {
			return nil, errors.BadRequest("invalid template id", "")
		}

		_, err = t.data.Mongo.Task.UpdateByIDIfExists(ctx,
			params.Id,
			mgz.Op().
				Set("status", "generating").
				Set("template", template),
		)

		if err != nil {
			log.Errorw("UpdateOneXXById err", err)
			return nil, err
		}

		var taskSegments []interface{}
		for _, x := range template.GetSegments() {

			x.Root = &projpb.Resource{
				XId:       template.XId,
				Url:       template.Url,
				CoverUrl:  template.CoverUrl,
				Commodity: template.Commodity,
			}

			taskSegments = append(taskSegments, &projpb.TaskSegment{
				Segment: x,
				TaskId:  task.XId,
				Task:    task,
				Status:  "textGenerating",
			})
		}

		_, err = t.data.Mongo.TaskSegment.GetCoreCollection().InsertMany(ctx, taskSegments)
		if err != nil {
			return nil, err
		}

		return task, nil

	case "updateStatus":

		_, err := t.data.Mongo.Task.UpdateByIDIfExists(ctx,
			params.Id,
			mgz.Op().
				Set("status", params.Params["status"]),
		)

		if err != nil {
			log.Errorw("UpdateOneXXById err", err)
			return nil, err
		}

	case "updateTargetRelatedItem":

		idx := params.Params["index"]

		_, err := t.data.Mongo.Task.UpdateByIDIfExists(ctx,
			params.Id,
			mgz.Op().
				Set(fmt.Sprintf("relatedItems.$[].selected"), false),
		)

		if err != nil {
			log.Errorw("UpdateOneXXById err", err)
			return nil, err
		}

		_, err = t.data.Mongo.Task.UpdateByIDIfExists(ctx,
			params.Id,
			mgz.Op().
				Set(fmt.Sprintf("relatedItems.%s.selected", idx), true).
				Set("status", "prepared"),
		)

		if err != nil {
			log.Errorw("UpdateOneXXById err", err)
			return nil, err
		}

	case "delete":
		err := t.data.Mongo.Task.DeleteByID(ctx, params.Id)
		if err != nil {
			log.Errorw("DestroyById err", err)
			return nil, err
		}
	}

	return nil, nil
}
