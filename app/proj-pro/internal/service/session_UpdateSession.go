package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"

	"github.com/go-kratos/kratos/v2/log"
)

func (t *SessionService) UpdateSession(ctx context.Context, params *projpb.UpdateSessionRequest) (*projpb.Session, error) {
	switch params.Action {
	case "startRemix":
		_, err := t.data.Mongo.Session.UpdateByIDIfExists(ctx,
			params.Id,
			mgz.Op().
				Set("status", "remixing"),
		)

		if err != nil {
			log.Errorw("UpdateOneXXById err", err)
			return nil, err
		}

	case "startSelectTemplate":
		_, err := t.data.Mongo.Session.UpdateByIDIfExists(ctx,
			params.Id,
			mgz.Op().
				Set("status", "templateSelecting"),
		)

		if err != nil {
			log.Errorw("UpdateOneXXById err", err)
			return nil, err
		}

	case "confirmSelectTemplate":

		templateId := params.Params["templateId"]

		template, err := t.data.Mongo.Template.GetById(ctx, templateId)
		if err != nil {
			return nil, err
		}

		_, err = t.data.Mongo.Session.UpdateByIDIfExists(ctx,
			params.Id,
			mgz.Op().
				Set("status", "generating").
				Set("template", template),
		)

		if err != nil {
			log.Errorw("UpdateOneXXById err", err)
			return nil, err
		}

		var sessionSegments []*projpb.SessionSegment
		for _, x := range template.Segments {

			x.Root = &projpb.Resource{
				XId:       template.XId,
				Url:       template.Url,
				CoverUrl:  template.CoverUrl,
				Commodity: template.Commodity,
			}

			sessionSegments = append(sessionSegments, &projpb.SessionSegment{
				Segment:   x,
				SessionId: params.Id,
				Status:    "textGenerating",
			})
		}

		_, err = t.data.Mongo.SessionSegment.InsertMany(ctx, sessionSegments...)
		if err != nil {
			return nil, err
		}

		//case "updateMode":
		//
		//	list, err := t.data.Mongo.SessionSegment.List(ctx, bson.M{"sessionId": params.Id})
		//	if err != nil {
		//		return nil, err
		//	}
		//
		//	for _, x := range list {
		//		if x.AssetId != "" {
		//			fmt.Println(x.AssetId)
		//			return nil, errors.Forbidden("", "")
		//		}
		//	}
		//
		//	_, err = t.data.Mongo.Task.UpdateByIDIfExists(ctx, params.Id,
		//		mongoz.Op().
		//			Set("mode", params.Params["mode"]),
		//	)
		//	if err != nil {
		//		log.Error(err)
		//		return nil, err
		//	}
		//
		//case "updateTargetChance":
		//
		//	task, err := t.data.Mongo.Task.GetById(ctx, params.Id)
		//	if err != nil {
		//		return nil, err
		//	}
		//
		//	chance := task.Commodity.Chances[int(params.ChanceIndex)]
		//
		//	_, err = t.data.Mongo.Task.UpdateByIDIfExists(ctx,
		//		params.Id,
		//		mongoz.Op().
		//			Set("status", "templateSelecting").
		//			Set("targetChance", chance),
		//	)
		//
		//	if err != nil {
		//		log.Errorw("UpdateOneXXById err", err)
		//		return nil, err
		//	}
		//case "updateTemplate":
		//
		//	task, err := t.data.Mongo.Task.GetById(ctx, params.Id)
		//	if err != nil {
		//		return nil, err
		//	}
		//
		//	template, err := t.data.GrpcClients.ProjAdminClient.GetTemplate(ctx, &projpb.GetGetTemplateRequest{
		//		Id: params.TemplateId,
		//	})
		//	if err != nil {
		//		return nil, err
		//	}
		//
		//	if template == nil {
		//		return nil, errors.BadRequest("invalid template id", "")
		//	}
		//
		//	_, err = t.data.Mongo.Task.UpdateByIDIfExists(ctx,
		//		params.Id,
		//		mongoz.Op().
		//			Set("status", "generating").
		//			Set("template", template),
		//	)
		//
		//	if err != nil {
		//		log.Errorw("UpdateOneXXById err", err)
		//		return nil, err
		//	}
		//
		//	var taskSegments []interface{}
		//	for _, x := range template.GetSegments() {
		//
		//		x.Root = &projpb.Resource{
		//			XId:       template.XId,
		//			Url:       template.Url,
		//			CoverUrl:  template.CoverUrl,
		//			Commodity: template.Commodity,
		//		}
		//
		//		taskSegments = append(taskSegments, &projpb.TaskSegment{
		//			Segment: x,
		//			TaskId:  task.XId,
		//			Task:    task,
		//			Status:  "textGenerating",
		//		})
		//	}
		//
		//	_, err = t.data.Mongo.TaskSegment.GetCoreCollection().InsertMany(ctx, taskSegments)
		//	if err != nil {
		//		return nil, err
		//	}
		//
		//	return task, nil
		//
		//case "updateStatus":
		//
		//	_, err := t.data.Mongo.Task.UpdateByIDIfExists(ctx,
		//		params.Id,
		//		mongoz.Op().
		//			Set("status", params.Params["status"]),
		//	)
		//
		//	if err != nil {
		//		log.Errorw("UpdateOneXXById err", err)
		//		return nil, err
		//	}
		//
		//case "updateTargetRelatedItem":
		//
		//	idx := params.Params["index"]
		//
		//	_, err := t.data.Mongo.Task.UpdateByIDIfExists(ctx,
		//		params.Id,
		//		mongoz.Op().
		//			Set(fmt.Sprintf("relatedItems.$[].selected"), false),
		//	)
		//
		//	if err != nil {
		//		log.Errorw("UpdateOneXXById err", err)
		//		return nil, err
		//	}
		//
		//	_, err = t.data.Mongo.Task.UpdateByIDIfExists(ctx,
		//		params.Id,
		//		mongoz.Op().
		//			Set(fmt.Sprintf("relatedItems.%s.selected", idx), true).
		//			Set("status", "prepared"),
		//	)
		//
		//	if err != nil {
		//		log.Errorw("UpdateOneXXById err", err)
		//		return nil, err
		//	}
		//
		//case "delete":
		//	err := t.data.Mongo.Task.DeleteByID(ctx, params.Id)
		//	if err != nil {
		//		log.Errorw("DestroyById err", err)
		//		return nil, err
		//	}
	}

	return nil, nil
}
