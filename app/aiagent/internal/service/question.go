package service

import (
	"context"
	errors2 "errors"
	"fmt"
	aiagentpb "store/api/aiagent"
	paymentpb "store/api/payment"
	"store/app/aiagent/configs"
	"store/pkg/krathelper"
	"store/pkg/middlewares/eventsource"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/wg"
	"strings"
	"time"

	"store/pkg/sdk/third/gemini"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/genai"
)

func (t TrackerService) AppendQuestion(ctx context.Context, params *aiagentpb.AppendQuestionParams) (*aiagentpb.Question, error) {

	userId := krathelper.RequireUserId(ctx)

	platform := krathelper.GetHeader(ctx, "Platform")

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "AppendQuestion",
		"sessionId", params.SessionId,
		//"scene", params.Scene,
		"prompt", params.Prompt,
		"platform", platform,
	))
	// todo 数据校验
	// todo 内存 获取session数据
	session, err := t.Data.Mongo.Session.GetById(ctx, params.SessionId)
	if err != nil {
		logger.Errorw("get session err", err)
		return nil, errors.BadRequest("invalidSession", "")
	}

	if len(session.Resources) == 0 {
		logger.Errorw("no resources found in session", err)
		return nil, errors.BadRequest("emptyResource", "")
	}

	//resources := aiagentpb.Resources(session.Resources)
	//// 将用户个人账号信息补充到resources
	//if params.Profile != nil {
	//	resources = append(resources, &aiagentpb.ResourceV2{
	//		SessionId: params.SessionId,
	//		Category:  "personalProfile",
	//		Content:   conv.S2J(params.Profile),
	//		MimeType:  "text/plain",
	//	})
	//}

	//cost := configs.GetCost("", params.Prompt.GetId())
	cost := int64(0)

	// 落库
	q, err := t.Data.Mongo.Question.InsertIfNX(ctx,
		bson.M{
			"session._id": params.SessionId,
			"prompt.id":   params.Prompt.GetId(),
			"status":      bson.M{"$in": []string{"created", "prepared", "generating"}},
		},
		&aiagentpb.Question{
			XId:       primitive.NewObjectID().Hex(),
			Session:   session,
			Prompt:    params.Prompt,
			UserId:    userId,
			CreatedAt: time.Now().Unix(),
			Cost: &aiagentpb.Cost{
				//Id:        credit.CostId,
				Amount: cost,
				//Remaining: credit.Remaining,
			},
			Status: "created",
			Extra: &aiagentpb.Question_Extra{
				Platform: platform,
			},
		})

	if err != nil {
		return nil, err
	}

	// 校验Quota
	creditState, err := t.Data.GrpcClients.PaymentClient.GetCreditState(ctx, &paymentpb.GetCreditStateParams{
		UserId: userId,
	})

	if err != nil {
		logger.Errorw("CostCredit err", err, "userId", userId)
		return nil, err
	}

	if creditState.Total-creditState.Used < cost {
		return nil, errors.BadRequest("exceeded", "")
	}

	_, err = t.Data.Mongo.Question.UpdateFieldsById(ctx, q.XId, bson.M{"status": "prepared"})
	if err != nil {
		return nil, err
	}

	return q, nil
}

func (t TrackerService) RetryQuestion(ctx context.Context, params *aiagentpb.RetryQuestionParams) (*aiagentpb.Question, error) {

	userId := krathelper.RequireUserId(ctx)

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "RetryQuestion",
		"questionId", params.QuestionId,
	))

	question, err := t.Data.Mongo.Question.GetById(ctx, params.QuestionId)
	if err != nil {
		return nil, err
	}

	_, err = t.Data.Mongo.Question.UpdateFieldsById(ctx, params.QuestionId, bson.M{"status": "created"})
	if err != nil {
		return nil, err
	}

	cost := configs.GetCost(question.Session.GetScene(), question.Prompt.GetId())

	// 校验Quota
	creditState, err := t.Data.GrpcClients.PaymentClient.GetCreditState(ctx, &paymentpb.GetCreditStateParams{
		UserId: userId,
	})

	if err != nil {
		logger.Errorw("CostCredit err", err, "userId", userId)
		return nil, err
	}

	if creditState.Total-creditState.Used < cost {
		return nil, errors.BadRequest("exceeded", "")
	}

	_, err = t.Data.Mongo.Question.UpdateFieldsById(ctx, question.XId, bson.M{"status": "prepared"})
	if err != nil {
		return nil, err
	}

	return question, nil
}

func (t TrackerService) SubmitQuestion(ctx context.Context, params *aiagentpb.SubmitQuestionParams) (*aiagentpb.Question, error) {

	userId := krathelper.RequireUserId(ctx)

	platform := krathelper.GetHeader(ctx, "Platform")

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "SubmitQuestion",
		"sessionId", params.SessionId,
		"scene", params.Scene,
		"questionId", params.QuestionId,
		"prompt", params.Prompt,
		"platform", platform,
	))
	// todo 数据校验
	// todo 内存 获取session数据
	session, err := t.Data.Mongo.Session.GetById(ctx, params.SessionId)
	if err != nil {
		logger.Errorw("get session err", err)
		return nil, errors.BadRequest("invalidSession", "")
	}

	if len(session.Resources) == 0 {
		logger.Errorw("no resources found in session", err)
		return nil, errors.BadRequest("emptyResource", "")
	}

	resources := aiagentpb.Resources(session.Resources)
	// 将用户个人账号信息补充到resources
	if params.Profile != nil {
		resources = append(resources, &aiagentpb.ResourceV2{
			SessionId: params.SessionId,
			Category:  "personalProfile",
			Content:   conv.S2J(params.Profile),
			MimeType:  "text/plain",
		})
	}

	//// 账号分析需要 补充 itemSummaries 资源
	//if params.Prompt.GetId() == "accountAnalysis" {
	//
	//	personalProfiles := helper.Filter(resources, func(res *aiagentpb.ResourceV2) bool {
	//		return res.Category == "personalProfile"
	//	})
	//
	//	logger.Debugw("personalProfiles", personalProfiles)
	//
	//	if len(personalProfiles) > 0 {
	//		pp := personalProfiles[0]
	//
	//		var acc *aiagentpb.Account
	//		_ = json.Unmarshal([]byte(pp.Content), &acc)
	//
	//		if acc != nil {
	//			var itemSummaries string
	//			if acc.Platform == "douyin" {
	//			} else {
	//				notes, err := t.Data.Tikhub.XhsGetUserNotes(ctx, acc.Extra["id"])
	//				if err != nil {
	//					logger.Errorw("get notes err", err, "acc", acc)
	//				}
	//
	//				itemSummaries = conv.S2J(notes)
	//			}
	//
	//			resources = append(resources, &aiagentpb.ResourceV2{
	//				SessionId: params.SessionId,
	//				Category:  "itemSummaries",
	//				Content:   itemSummaries,
	//				MimeType:  "text/plain",
	//			})
	//		}
	//
	//	}
	//
	//}

	cost := configs.GetCost(params.Scene, params.Prompt.GetId())

	// 落库
	q, err := t.Data.Mongo.Question.InsertIfNX(ctx,
		bson.M{
			"session._id": params.SessionId,
			"prompt.id":   params.Prompt.GetId(),
			"status":      bson.M{"$in": []string{"created", "prepared", "generating"}},
		},
		&aiagentpb.Question{
			XId:       primitive.NewObjectID().Hex(),
			Session:   session,
			Prompt:    params.Prompt,
			UserId:    userId,
			CreatedAt: time.Now().Unix(),
			Cost: &aiagentpb.Cost{
				//Id:        credit.CostId,
				Amount: cost,
				//Remaining: credit.Remaining,
			},
			Status: "created",
			Extra: &aiagentpb.Question_Extra{
				Platform: platform,
			},
		})

	if err != nil {
		return nil, err
	}

	// 校验Quota
	creditState, err := t.Data.GrpcClients.PaymentClient.GetCreditState(ctx, &paymentpb.GetCreditStateParams{
		UserId: userId,
	})

	if err != nil {
		logger.Errorw("CostCredit err", err, "userId", userId)
		return nil, err
	}

	if creditState.Total-creditState.Used < cost {
		return nil, errors.BadRequest("exceeded", "")
	}

	_, err = t.Data.Mongo.Question.UpdateFieldsById(ctx, q.XId, bson.M{"status": "prepared"})
	if err != nil {
		return nil, err
	}

	return q, nil
}

func (t TrackerService) UpdateQuestionStatus(ctx context.Context, params *empty.Empty) (*empty.Empty, error) {

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "UpdateQuestionStatus",
	))

	//logger.Debugw("update question status", "ing")

	questions, err := t.Data.Mongo.Question.List(ctx, bson.M{
		"status":    bson.M{"$in": []string{"generating"}},
		"createdAt": bson.M{"$lte": time.Now().Add(-5 * time.Minute).Unix()}, // 5分钟没生成直接失败 进入下一次任务
	})

	ids := helper.Mapping(questions, func(x *aiagentpb.Question) string {
		return x.XId
	})

	_, err = t.Data.Mongo.Question.UpdateFields(ctx, bson.M{"_id": bson.M{"$in": ids}}, bson.M{"status": "failed", "reason": "time out"})
	if err != nil {
		logger.Errorw("UpdateQuestionStatus err", err)
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (t TrackerService) GenerateAnswerChunks(ctx context.Context, params *empty.Empty) (*empty.Empty, error) {

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "GenerateAnswerChunks",
	))

	questions, err := t.Data.Mongo.Question.List(ctx, bson.M{"status": bson.M{"$in": []string{"prepared", "toRetry"}}})
	if err != nil {
		logger.Errorw("list questions err", err)
		return nil, err
	}

	//logger.Debugw("list questions for GenerateAnswerChunks", len(questions))

	if len(questions) == 0 {
		return &empty.Empty{}, nil
	}

	wg.WaitGroup(context.Background(), questions, t.generateAnswer)

	return &empty.Empty{}, nil
}

func (t TrackerService) generateAnswer(ctx context.Context, question *aiagentpb.Question) error {

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "generateAnswer",
		"sessionId", question.Session.GetXId(),
		"questionId", question.XId,
	))

	_, err := t.Data.Mongo.Question.UpdateFieldsById(ctx, question.XId, bson.M{"status": "generating"})
	if err != nil {
		return err
	}

	// 包装Prompt
	scene := question.Session.GetScene()
	prompt := question.Prompt
	//resources := question.Session.GetResources()

	genaiClient := t.Data.GenaiFactory.Get()

	parts, err := t.packSystemPromptV3(ctx, question)
	if err != nil {
		_, err = t.Data.Mongo.Question.UpdateFieldsById(ctx, question.XId, bson.M{
			"status": "prepared", "reason": "upload failed:" + err.Error(),
		})

		return err
	}

	//model := service.Data.GenaiClient.GenerativeModel("gemini-2.0-flash")
	//model := t.Data.GenaiClient.GenerativeModel("gemini-2.5-pro-exp-03-25")
	//model := genaiClient.GenerativeModel("gemini-2.5-pro-preview-05-06")

	modelName := configs.GetModel(scene, prompt.GetId())

	// 重试过程中更换模型
	if question.GetState().GetRetryTimes() > 3 {
		modelName = configs.GetModalPro()
	}

	key := fmt.Sprintf("chunks:%s", question.XId)

	onMessage := func(index int, text string) error {

		// 可能是重试或者刚开始 重置redis
		if index == 0 {
			t.Data.Redis.Del(ctx, key)
		}

		// 防止被去重, 获取的时候要删掉
		text = conv.Str(index) + text
		t.Data.Redis.ZAdd(ctx, key, redis.Z{Score: float64(index), Member: text})
		// 每次延长过期时间
		t.Data.Redis.ExpireAt(ctx, key, time.Now().Add(5*time.Hour))

		return nil
	}

	//retryTimes, err := helper.Retry(ctx, func(ctx context.Context) error {
	//}, 3)

	logger.Debugw("pre doGenerate log", "", "parts", parts)

	err = t.doGenerate(ctx, genaiClient, modelName, parts, onMessage)
	if err != nil {

		//重试超过 最大次就不再重试了 直接失败
		if question.State.GetRetryTimes() > 5 {

			_, err = t.Data.Mongo.Question.UpdateFieldsById(ctx, question.XId, bson.M{
				"status": "failed", "reason": "exceeded retry times:" + err.Error(), "model": modelName,
			})

			return err
		}

		_, err = t.Data.Mongo.Question.UpdateFieldsById(ctx, question.XId, bson.M{
			"status": "toRetry", "reason": err.Error(), "model": modelName,
			"state": bson.M{"retryTimes": question.State.GetRetryTimes() + 1},
		})

		return err
	}

	chunks, err := t.getAnswerChunks(ctx, question.XId, "")
	if err != nil {
		return err
	}

	answerText := chunks.FullText()
	if answerText == "" {
		_, err = t.Data.Mongo.Question.UpdateFieldsById(ctx, question.XId, bson.M{
			"status": "toRetry", "reason": "empty text", "model": modelName,
			"state": bson.M{"retryTimes": question.State.GetRetryTimes() + 1},
		})
		return errors2.New("empty text")
	}
	// 更新question状态
	_, err = t.Data.Mongo.Question.UpdateFieldsById(ctx, question.XId, bson.M{
		"status": "completed",
		"model":  modelName,
		"answer": &aiagentpb.Answer{
			Text: chunks.FullText(),
		},
	})
	if err != nil {
		return err
	}

	_, err = t.Data.GrpcClients.PaymentClient.CostCredit(ctx, &paymentpb.CheckCreditParams{
		UserId: question.UserId,
		Cost:   question.Cost.GetAmount(),
		Key:    question.XId,
	})
	if err != nil {
		logger.Errorw("cost credit err", err, "question", question)
		return err
	}

	return nil
}

func (t TrackerService) doGenerate(ctx context.Context, genaiClient *gemini.Client, modelName string, parts []*genai.Part, onMessage func(index int, text string) error) error {

	var index int

	for resp, err := range genaiClient.GenerateContentStream(ctx, gemini.GenerateContentRequest{
		Model: modelName,
		Parts: parts,
	}) {
		if err != nil {
			return err
		}

		part := gemini.ResponseToString(resp)

		if strings.TrimSpace(part) == "" {
			continue
		}

		err = onMessage(index, part)
		if err != nil {
			return err
		}

		index++
	}

	return onMessage(9999, "")
}

func (t TrackerService) GetAnswerChunks(ctx context.Context, params *aiagentpb.GetAnswerChunksParams) (*aiagentpb.AnswerChunks, error) {
	return t.getAnswerChunks(ctx, params.QuestionId, params.StartChunkId)
}

func (t TrackerService) getAnswerChunks(ctx context.Context, questionId, startChunkId string) (*aiagentpb.AnswerChunks, error) {

	key := fmt.Sprintf("chunks:%s", questionId)

	zs, err := t.Data.Redis.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min: startChunkId,
		Max: "+inf",
	}).Result()

	if err != nil {
		return nil, err
	}

	if len(zs) == 0 {
		return &aiagentpb.AnswerChunks{}, nil
	}

	var chunks []*aiagentpb.AnswerChunk

	var done bool
	for _, z := range zs {

		if z.Score == 9999 {
			done = true
			continue
		}

		score := conv.Str(z.Score)
		text := strings.Replace(z.Member.(string), score, "", 1)

		chunks = append(chunks, &aiagentpb.AnswerChunk{
			ChunkId: score,
			Text:    text,
		})
	}

	return &aiagentpb.AnswerChunks{
		Chunks:       chunks,
		StartChunkId: conv.Str(zs[0].Score),
		EndChunkId:   conv.Str(zs[len(zs)-1].Score),
		Done:         done,
	}, nil
}

func (t TrackerService) GetAnswerChunksStream(ctx context.Context, params *aiagentpb.GetAnswerChunksStreamParams) (int, error) {
	esCtx := ctx.(eventsource.Ctx)
	//userId := esCtx.UserId()

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "GetAnswerChunksStream",
		"questionId", params.QuestionId,
	))

	lastChunkId := ""

	for {
		time.Sleep(500 * time.Millisecond)

		chunks, err := t.getAnswerChunks(ctx, params.QuestionId, lastChunkId)
		if err != nil {
			return 0, err
		}

		if len(chunks.Chunks) > 0 {
			logger.Debugw("got chunks ", "",
				"lastChunkId", lastChunkId,
				"StartChunkId", chunks.StartChunkId,
				"EndChunkId", chunks.EndChunkId)
		}

		for _, x := range chunks.Chunks {
			// 去重 期望redis的结果是左开，但是实际结果是左闭，需要去掉
			if x.ChunkId == lastChunkId {
				continue
			}

			esCtx.WriteV2(x.ChunkId, x.Text)
		}

		if chunks.Done {
			esCtx.WriteDone("")
			break
		}

		lastChunkId = chunks.EndChunkId
	}

	return 0, nil
}
