package service

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	aiagentpb "store/api/aiagent"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"time"
)

func (t TrackerService) ListSessions(ctx context.Context, params *aiagentpb.ListSessionsParams) (*aiagentpb.SessionList, error) {
	userId := krathelper.FindUserId(ctx)

	if userId == "" {
		return &aiagentpb.SessionList{}, nil
	}

	filters := bson.M{
		"userId": userId,
		//"scene":  params.Scene,
		"status": bson.M{
			"$ne": "created",
		},
	}

	if params.Scene != "" {
		filters["scene"] = params.Scene
	}

	if params.StartTs > 0 {
		filters["createdAt"] = bson.M{"$gte": params.StartTs}
	}

	if params.EndTs > 0 {
		filters["createdAt"] = bson.M{"$lte": params.EndTs}
	}

	questions, _, err := t.Data.Mongo.Session.ListAndCount(ctx,
		filters,
		mgz.
			Paging(params.Page, params.Size).
			SetSort(bson.M{"createdAt": -1}),
	)
	if err != nil {
		return nil, err
	}

	return &aiagentpb.SessionList{
		List: aiagentpb.Sessions(questions).Safe(),
	}, nil
}

func (t TrackerService) CreateQuickSession(ctx context.Context, params *aiagentpb.CreateQuickSessionParams) (*aiagentpb.Session, error) {

	userId := krathelper.RequireUserId(ctx)

	item, err := t.Data.Mongo.Item.GetById(ctx, params.ItemId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// 如何已经存在直接返回

	olds, err := t.Data.Mongo.Session.List(ctx, bson.M{
		"lookupKey": item.XId,
		"userId":    userId,
	})
	if err != nil {
		return nil, err
	}

	if len(olds) > 0 {
		_, _ = t.Data.Mongo.Session.UpdateFieldsById(ctx,
			olds[0].XId,
			bson.M{"createdAt": time.Now().Unix()},
		)

		return olds[0], nil
	}

	promptId := "analysis"
	scene := "analysis"

	report := item.Reports[promptId]
	if report == "" {
		log.Errorw("err", "no analysis report", "itemId", params.ItemId)
		return nil, nil
	}

	// 创建session
	sessionId := primitive.NewObjectID().Hex()

	newSession := &aiagentpb.Session{
		XId:       sessionId,
		UserId:    userId,
		CreatedAt: time.Now().Unix(),
		Scene:     scene,
		LookupKey: item.XId,
		Resources: []*aiagentpb.ResourceV2{
			{
				XId:          item.XId,
				SessionId:    sessionId,
				Url:          item.Url,
				Profile:      item.Profile,
				Title:        item.Title,
				Desc:         item.Desc,
				UploadUserId: userId,
				CreatedAt:    time.Now().Unix(),
				MimeType:     "video/mp4",
				CoverUrl:     item.Cover,
				//GenaiUri:     item.GenaiUri,
			},
		},
	}

	_, err = t.Data.Mongo.Session.InsertOne(ctx, newSession)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	questionId := primitive.NewObjectID().Hex()

	// 创建 answer
	answerId := primitive.NewObjectID().Hex()
	newAnswer := &aiagentpb.Answer{
		XId: answerId,
		Question: &aiagentpb.Question{
			XId: questionId,
		},
		Text:      report,
		UserId:    userId,
		CreatedAt: time.Now().Unix(),
	}
	_, err = t.Data.Mongo.Answer.InsertOne(ctx, newAnswer)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// 创建question
	_, err = t.Data.Mongo.Question.InsertOne(ctx, &aiagentpb.Question{
		XId:       questionId,
		Session:   newSession,
		Prompt:    &aiagentpb.Prompt{Id: promptId},
		UserId:    userId,
		CreatedAt: time.Now().Unix(),
		Cost:      &aiagentpb.Cost{
			//Id:        credit.CostId,
			//Amount:    credit.Amount,
			//Remaining: credit.Remaining,
		},
		Answers: []*aiagentpb.Answer{newAnswer},
		Answer:  newAnswer,
		Status:  "completed",
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return newSession, nil
}

func (t TrackerService) GetSession(ctx context.Context, params *aiagentpb.GetSessionParams) (*aiagentpb.Session, error) {

	//userId := krathelper.RequireUserId(ctx)

	res, err := t.Data.Mongo.Session.GetById(ctx, params.Id)
	if err != nil {
		log.Error("Session err", err, "params", params)
		return nil, errors.BadRequest("no session found", "")
	}

	return res.Safe(), nil
}

func (t TrackerService) CreateSessionV3(ctx context.Context, params *aiagentpb.CreateSessionV3Params) (*aiagentpb.Session, error) {

	userId := krathelper.RequireUserId(ctx)

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "CreateSessionV3",
		"sessionId", params.SessionId,
		"scene", params.Scene,
	))
	// 资源

	var resources []*aiagentpb.ResourceV2

	for _, x := range params.Resources {

		resourceId := primitive.NewObjectID().Hex()

		newResource := &aiagentpb.ResourceV2{
			XId:          resourceId,
			Url:          x.GetUrl(),
			UploadUserId: userId,
			CreatedAt:    time.Now().Unix(),
			MimeType:     x.GetMimeType(),
			SessionId:    params.SessionId,
			Name:         x.GetName(),
			Category:     x.GetCategory(),
			Content:      x.Content,
			Meta:         x.GetMeta(),
			Uri:          x.GetUri(),
			CoverUrl:     x.GetCoverUrl(),
		}

		resources = append(resources, newResource)
	}

	if params.Resource != nil {
		resourceId := primitive.NewObjectID().Hex()

		newResource := &aiagentpb.ResourceV2{
			XId:          resourceId,
			Url:          params.Resource.GetUrl(),
			UploadUserId: userId,
			CreatedAt:    time.Now().Unix(),
			MimeType:     params.Resource.GetMimeType(),
			SessionId:    params.SessionId,
			Name:         params.Resource.GetName(),
			Category:     params.Resource.GetCategory(),
			Content:      params.Resource.GetContent(),
			Meta:         params.Resource.GetMeta(),
			Uri:          params.Resource.GetUri(),
			CoverUrl:     params.Resource.GetCoverUrl(),
		}

		resources = append(resources, newResource)
	}

	// 账号分析需要 补充 itemSummaries 资源
	if params.Scene == "accountAnalysis" {

		personalProfiles := helper.Filter(resources, func(res *aiagentpb.ResourceV2) bool {
			return res.Category == "personalProfile"
		})

		logger.Debugw("personalProfiles", personalProfiles)

		if len(personalProfiles) > 0 {
			pp := personalProfiles[0]

			var acc *aiagentpb.Account
			_ = json.Unmarshal([]byte(pp.Content), &acc)

			if acc != nil {
				var itemSummaries string
				if acc.Platform == "douyin" {
				} else {
					notes, err := t.Data.Tikhub.XhsGetUserNotes(ctx, acc.Extra["id"])
					if err != nil {
						logger.Errorw("get notes err", err, "acc", acc)
					}

					itemSummaries = conv.S2J(notes)
				}

				resources = append(resources, &aiagentpb.ResourceV2{
					SessionId: params.SessionId,
					Category:  "itemSummaries",
					Content:   itemSummaries,
					MimeType:  "text/plain",
				})
			}

		}

	}

	err := t.Data.Mongo.Resource.InsertBatch(ctx, resources)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	newSession := &aiagentpb.Session{
		XId:       params.SessionId,
		UserId:    userId,
		CreatedAt: time.Now().Unix(),
		Scene:     params.Scene,
		Resources: resources,
		Version:   params.Version,
	}

	_, err = t.Data.Mongo.Session.InsertIfNotExistsByID(ctx, params.SessionId, newSession)
	if err != nil {
		log.Error("InsertOne err", err, "newSession", newSession)
		return nil, err
	}

	return newSession.Safe(), nil
}

func (t TrackerService) CreateSessionV2(ctx context.Context, params *aiagentpb.CreateSessionV2Params) (*aiagentpb.Session, error) {
	userId := krathelper.RequireUserId(ctx)

	newSession := &aiagentpb.Session{
		XId:       params.SessionId,
		UserId:    userId,
		CreatedAt: time.Now().Unix(),
		Scene:     params.Scene,
		Version:   params.Version,
	}

	_, err := t.Data.Mongo.Session.InsertIfNotExistsByID(ctx, params.SessionId, newSession)
	if err != nil {
		log.Error("InsertOne err", err, "newSession", newSession)
		return nil, err
	}

	return newSession.Safe(), nil

}

func (t TrackerService) CreateSession(ctx context.Context, params *aiagentpb.CreateSessionParams) (*aiagentpb.Session, error) {
	userId := krathelper.RequireUserId(ctx)

	res, err := t.Data.Mongo.Resource.GetById(ctx, params.ResourceId)
	if err != nil {
		log.Error("Resource.GetById err", err, "params.ResourceId", params.ResourceId)
		return nil, errors.BadRequest("no resource found", "")
	}

	newSession := &aiagentpb.Session{
		XId:       primitive.NewObjectID().Hex(),
		Resource:  res,
		UserId:    userId,
		CreatedAt: time.Now().Unix(),
		Scene:     params.Scene,
	}

	_, err = t.Data.Mongo.Session.InsertOne(ctx, newSession)
	if err != nil {
		log.Error("InsertOne err", err, "newSession", newSession)
		return nil, err
	}

	return newSession.Safe(), nil
}
