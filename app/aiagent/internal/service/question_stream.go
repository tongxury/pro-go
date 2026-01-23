package service

import (
	"context"
	"fmt"
	trackerpb "store/api/aiagent"
	paymentpb "store/api/payment"
	userpb "store/api/user"
	"store/app/aiagent/configs"
	"store/pkg/krathelper"
	"store/pkg/middlewares/eventsource"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/third/gemini"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/genai"
)

func (t TrackerService) ListQuestions(ctx context.Context, params *trackerpb.ListQuestionsParams) (*trackerpb.QuestionList, error) {
	userId := krathelper.FindUserId(ctx)

	if userId == "" {
		return &trackerpb.QuestionList{}, nil
	}

	filter := bson.M{
		"userId": userId,
	}

	if params.SessionId != "" {
		filter["session._id"] = params.SessionId
	}

	if len(params.GetStatus()) > 0 {
		filter["status"] = bson.M{"$in": params.GetStatus()}
	}

	if params.Ongoing {
		filter["status"] = bson.M{"$in": []string{"prepared", "generating", "toRetry", "failed"}}
	}

	if params.StartTs > 0 {
		filter["createdAt"] = bson.M{"$gte": params.StartTs}
	}

	opts := options.Find().SetSort(bson.M{"createdAt": 1})
	if params.Sort != "" {
		parts := strings.Split(params.Sort, "_")
		if len(parts) == 2 {
			opts.SetSort(bson.M{parts[0]: conv.Int64(parts[1])})
		}
	}

	questions, _, err := t.Data.Mongo.Question.ListAndCount(ctx,
		filter,
		opts,
	)
	if err != nil {
		return nil, err
	}

	return &trackerpb.QuestionList{
		List: trackerpb.Questions(questions).Safe(),
	}, nil
}

//func (t TrackerService) uploadToGenai(ctx context.Context, client *genai.Client, url, mimeType string) (string, error) {
//	result, err := resty.New().R().Get(url)
//	if err != nil {
//		return "", err
//	}
//
//	fileBody := bytes.NewReader(result.Body())
//
//	t1 := time.Now()
//	log.Debugw("start to upload file to genai time", t1)
//
//	ctx = context.Background()
//
//	opts := genai.UploadFileOptions{DisplayName: "", MIMEType: mimeType}
//	response, err := client.UploadFile(ctx, "", fileBody, &opts)
//	if err != nil {
//		log.Errorw("upload file error", err)
//		return "", err
//	}
//
//	// 校验
//	for response.State == genai.FileStateProcessing {
//		fmt.Print(".")
//		// Sleep for 10 seconds
//		time.Sleep(1 * time.Second)
//		// Fetch the file from the API again.
//		response, err = client.GetFile(ctx, response.Name)
//		if err != nil {
//			log.Fatal(err)
//			continue
//		}
//
//		log.Debugw("GenaiClient.GetFile response", response)
//	}
//
//	return response.URI, nil
//}

func (t TrackerService) packSystemPromptV2(ctx context.Context, genaiClient *gemini.Client, scene string, prompt *trackerpb.Prompt, resources trackerpb.Resources) ([]*genai.Part, error) {

	//log.Debugw("start to pack system prompt v2", "", "scene", scene, "prompt", prompt, "resources", resources, "resources.IsAllMediaImages() ", resources.IsAllMediaImages())

	configKey := fmt.Sprintf("%s", prompt.GetId())

	// todo
	if prompt.GetId() == "limitAnalysis" {
		personalProfile := resources.FindOneByCategory("personalProfile")
		if personalProfile != nil && personalProfile.Meta["platform"] == "douyin" {
			configKey += "_douyin"
		}
	}

	config, err := t.Data.Mongo.PromptConfig.Get(ctx, configKey)

	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, nil
	}

	var results []*genai.Part
	for _, x := range config.Parts {
		if x.Category == "" || x.Category == "text" {
			results = append(results, gemini.NewTextPart(x.Value))
		}

		if x.Category == "resource" {
			// 文件
			if x.Value == "" {

				for _, xx := range resources {
					if xx.Category != "" {
						continue
					}

					p, err := gemini.NewMediaPart(xx.Url, xx.MimeType)
					if err != nil {
						log.Errorw("NewMediaPart err", err, "url", xx.Url)
						return nil, err
					}

					results = append(results, p)
				}
			} else {
				// 文本
				y := resources.FindOneByCategory(x.Value)
				if y != nil {

					//content := strings.ReplaceAll(y.Content, "")

					results = append(results, gemini.NewTextPart(y.Content))
				}
			}

		}
	}

	if len(results) == 0 {
		return nil, errors.BadRequest("invalidPromptId", "")
	}

	return results, nil
}

// 根据场景包装 参数
func (t TrackerService) packSystemPrompt(ctx context.Context, scene string, prompt *trackerpb.Prompt, resources trackerpb.Resources) (string, error) {

	settings, err := t.Data.Mongo.Settings.Get(ctx)
	if err != nil {
		return "", err
	}

	var promptText string

	log.Debugw("start to pack system prompt", "", "scene", scene, "prompt", prompt, "resources", resources, "resources.IsAllMediaImages() ", resources.IsAllMediaImages())

	switch scene {
	// 爆款分析
	case "analysis":
		switch prompt.GetId() {
		case "analysisImages":
			promptText = settings.Prompt["analysis_analysisImages"]
		}
		if promptText == "" {
			// 图文
			if resources.IsAllMediaImages() {
				promptText = settings.Prompt["analysis_analysis_imageText"]
			} else {
				promptText = settings.Prompt["analysis_analysis_video"]
			}
		}

	case "coverAnalysis":
		if prompt.GetId() == "coverAnalysisImages" {
			promptText = settings.Prompt["coverAnalysis_coverAnalysisImages"]
		}

		if prompt.GetId() == "coverAnalysis" {
			promptText = settings.Prompt["coverAnalysis_coverAnalysis"]
		}
	case "limitAnalysis":
		switch prompt.GetId() {
		case "limitAnalysisImages":

			//configs, err := t.Data.Mongo.PromptConfig.Get(ctx, "limitAnalysis_limitAnalysisImages")
			//if err != nil {
			//	return "", err
			//}
			//
			//fmt.Println(configs)

			promptText = settings.Prompt["limitAnalysis_limitAnalysisImages"]
		}
		//if resources.IsAllMediaImages() {
		//	promptText = settings.Prompt["limitAnalysis_limitAnalysis_imageText"]
		//} else {
		//	//promptText = settings.Prompt["limitAnalysis_limitAnalysis_video"]
		//}
	case "preAnalysis":

		switch prompt.GetId() {
		case "preAnalysisImages":
			promptText = settings.Prompt["preAnalysis_preAnalysisImages"]
		}

		if resources.IsAllMediaImages() {
			promptText = settings.Prompt["preAnalysis_preAnalysis_imageText"]
		} else {
			//promptText = settings.Prompt["limitAnalysis_limitAnalysis_video"]
		}
	}

	if promptText == "" {
		return "", errors.BadRequest("invalidPromptId", "")
	}

	return promptText, nil
}

func (t TrackerService) CreateQuestionV3(ctx context.Context, params *trackerpb.CreateQuestionParams) (int, error) {
	esCtx := ctx.(eventsource.Ctx)

	userId := esCtx.UserId()

	if !params.Retry {
		if !helper.InSlice(params.GetPrompt().GetId(), []string{"extractText", "editingTechAnalysis", "effectAnalysis", "soundAnalysis"}) {
			list, err := t.Data.Mongo.Question.List(ctx, bson.M{
				"userId":      userId,
				"session._id": params.SessionId,
				"prompt.id":   params.Prompt.GetId(),
				//"status":      "completed",
			})
			if err != nil {
				return 0, err
			}

			if len(list) != 0 {
				log.Debugw("question already exists", "", "userId", userId, "sessionId", params.SessionId, "promptId", params.Prompt.GetId())
				return 0, errors.BadRequest("duplicate", "")
			}
		}
	}

	// todo 内存
	session, err := t.Data.Mongo.Session.GetById(ctx, params.SessionId)
	if err != nil {
		return 0, errors.BadRequest("invalidSession", "")
	}

	if len(session.Resources) == 0 {
		return 0, errors.BadRequest("emptyResource", "")
	}

	resources := trackerpb.Resources(session.Resources)
	// 将用户个人账号信息补充到resources
	if params.Profile != nil {
		resources = append(resources, &trackerpb.ResourceV2{
			SessionId: params.SessionId,
			Category:  "personalProfile",
			Content:   conv.S2J(params.Profile),
			MimeType:  "text/plain",
		})
	}

	// 包装Prompt
	systemContent, err := t.packSystemPrompt(ctx, params.Scene, params.Prompt, resources)

	if err != nil {
		return 0, err
	}

	if systemContent == "" {
		return 10400, errors.BadRequest("invalidPrompt", "")
	}
	params.Prompt.SystemContent = systemContent

	log.Debugw("params", params, "prompt", params.Prompt, "resources", resources)

	// 校验credit
	cost := configs.GetCost(params.Scene, params.Prompt.GetId())

	// 落库
	q, err := t.Data.Mongo.Question.InsertIfNX(ctx,
		bson.M{
			"session._id": params.SessionId,
			"prompt.id":   params.Prompt.GetId(),
			"status":      "created",
		},
		&trackerpb.Question{
			XId:       primitive.NewObjectID().Hex(),
			Session:   session,
			Prompt:    params.Prompt,
			UserId:    userId,
			CreatedAt: time.Now().Unix(),
			Cost: &trackerpb.Cost{
				//Id:        credit.CostId,
				Amount: cost,
				//Remaining: credit.Remaining,
			},
			Status: "created",
		})
	if err != nil {
		return 500, err
	}

	log.Debugw("InsertOne", "", "question.id", q.XId)

	credit, err := t.Data.GrpcClients.PaymentClient.CostCredit(ctx, &paymentpb.CheckCreditParams{
		UserId: userId,
		Cost:   cost,
		Key:    q.XId,
	})

	if err != nil {
		log.Errorw("CostCredit err", err, "userId", userId)
		return 500, err
	}
	if !credit.Ok {
		return 10501, errors.BadRequest("exceeded", "")
	}

	// 包装Genai参数
	var parts []*genai.Part

	geminiClient := t.Data.GenaiFactory.Get()

	// 资源
	for _, x := range resources {
		if strings.HasPrefix(x.MimeType, "video/") || strings.HasPrefix(x.MimeType, "image/") {

			p, err := gemini.NewMediaPart(x.Url, x.MimeType)
			if err != nil {
				log.Errorw("NewMediaPart err", err, "url", x.Url)
				return 0, err
			}

			parts = append(parts, p)
		}
		if strings.HasPrefix(x.MimeType, "text/") {
			parts = append(parts, gemini.NewTextPart(x.PromptText()))
		}

	}

	// prompt
	parts = append(parts, gemini.NewTextPart(systemContent))

	log.Debugw("genaiParams", "", "parts", parts)

	//parts := []genai.Part{
	//	genai.FileData{
	//		MIMEType: "video/mp4",
	//		URI:      session.Resource.GenaiUri,
	//	},
	//	genai.Text(params.Prompt.SystemContent),
	//}
	//

	_, err = t.Data.Mongo.Question.UpdateOne(ctx,
		bson.M{"_id": q.XId},
		bson.M{"$set": bson.M{"status": "generating"}},
	)
	if err != nil {

	}

	var answerText string

	for resp, err := range geminiClient.GenerateContentStream(ctx, gemini.GenerateContentRequest{
		Model: "gemini-2.5-flash-preview-05-20",
		Parts: parts,
	}) {
		if err != nil {
			//ERROR iterate fail err=googleapi: Error 403: You do not have permission to access the File xfr2ktwx3cb5 or it may not exist.

			log.Errorw("iterate fail err", err)

			return 500, err
		}

		part := gemini.ResponseToString(resp)

		if strings.TrimSpace(part) == "" {
			continue
		}

		log.Debugw("answer part", part)

		answerText += part
		esCtx.Write(part)
	}

	if answerText == "" {
		return 500, errors.BadRequest("invalidAnswer", "")
	}

	// answer落库
	answerId := primitive.NewObjectID().Hex()
	newAnswer := &trackerpb.Answer{
		XId: answerId,
		Question: &trackerpb.Question{
			XId: q.XId,
		},
		Text:      answerText,
		UserId:    userId,
		CreatedAt: time.Now().Unix(),
	}
	_, err = t.Data.Mongo.Answer.InsertOne(ctx, newAnswer)
	if err != nil {
		return 500, err
	}

	_, err = t.Data.Mongo.Question.UpdateOne(ctx,
		bson.M{"_id": q.XId},
		bson.M{
			"$push": bson.M{"answers": newAnswer},
			"$set": bson.M{
				"status": "completed",
			},
		},
		options.Update().SetUpsert(true))
	if err != nil {
		return 500, err
	}

	question, err := t.Data.Mongo.Question.GetById(ctx, q.XId)
	if err != nil {
		return 500, err
	}

	esCtx.WriteDone(question.Safe())
	return 0, nil
}

func (t TrackerService) CreateQuestionV5(ctx context.Context, params *trackerpb.CreateQuestionParams) (int, error) {
	esCtx := ctx.(eventsource.Ctx)

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "CreateQuestionV5",
		"sessionId", params.SessionId,
		"scene", params.Scene,
	))

	userId := esCtx.UserId()

	// 去重
	if !params.Retry {
		if !helper.InSlice(params.GetPrompt().GetId(), []string{"extractText", "editingTechAnalysis", "effectAnalysis", "soundAnalysis"}) {
			list, err := t.Data.Mongo.Question.List(ctx, bson.M{
				"userId":      userId,
				"session._id": params.SessionId,
				"prompt.id":   params.Prompt.GetId(),
				"status":      "completed",
			})
			if err != nil {
				return 0, err
			}

			if len(list) != 0 {
				logger.Debugw("question already exists", "", "userId", userId, "sessionId", params.SessionId, "promptId", params.Prompt.GetId())
				return 0, errors.BadRequest("duplicate", "")
			}
		}
	}

	// todo 内存 获取session数据
	session, err := t.Data.Mongo.Session.GetById(ctx, params.SessionId)
	if err != nil {
		logger.Errorw("get session err", err)
		return 0, errors.BadRequest("invalidSession", "")
	}

	if len(session.Resources) == 0 {
		logger.Errorw("no resources found in session", err)
		return 0, errors.BadRequest("emptyResource", "")
	}

	resources := trackerpb.Resources(session.Resources)
	// 将用户个人账号信息补充到resources
	if params.Profile != nil {
		resources = append(resources, &trackerpb.ResourceV2{
			SessionId: params.SessionId,
			Category:  "personalProfile",
			Content:   conv.S2J(params.Profile),
			MimeType:  "text/plain",
		})
	}

	cost := configs.GetCost(params.Scene, params.Prompt.GetId())

	// 落库
	q, err := t.Data.Mongo.Question.InsertIfNX(ctx,
		bson.M{
			"session._id": params.SessionId,
			"prompt.id":   params.Prompt.GetId(),
			"status":      "created",
		},
		&trackerpb.Question{
			XId:       primitive.NewObjectID().Hex(),
			Session:   session,
			Prompt:    params.Prompt,
			UserId:    userId,
			CreatedAt: time.Now().Unix(),
			Cost: &trackerpb.Cost{
				//Id:        credit.CostId,
				Amount: cost,
				//Remaining: credit.Remaining,
			},
			Status: "created",
		})
	if err != nil {
		return 500, err
	}

	logger.Debugw("InsertOne", "", "question.id", q.XId)

	// 校验credit
	credit, err := t.Data.GrpcClients.PaymentClient.CostCredit(ctx, &paymentpb.CheckCreditParams{
		UserId: userId,
		Cost:   cost,
		Key:    q.XId,
	})

	if err != nil {
		log.Errorw("CostCredit err", err, "userId", userId)
		return 500, err
	}
	if !credit.Ok {
		return 10501, errors.BadRequest("exceeded", "")
	}

	// 包装Prompt
	genaiClient := t.Data.GenaiFactory.Get()

	parts, err := t.packSystemPromptV2(ctx, genaiClient, params.Scene, params.Prompt, resources)

	if err != nil {
		return 0, err
	}

	//model := service.Data.GenaiClient.GenerativeModel("gemini-2.0-flash")
	//model := t.Data.GenaiClient.GenerativeModel("gemini-2.5-pro-exp-03-25")
	//model := genaiClient.GenerativeModel("gemini-2.5-pro-preview-05-06")

	modelName := configs.GetModel(params.Scene, params.Prompt.GetId())

	//model := genaiClient.GenerativeModel(modelName)

	logger.Debugw("params", params, "modelName", modelName, "cost", cost, "prompt", params.Prompt, "resources", len(resources), "parts", parts)

	_, err = t.Data.Mongo.Question.UpdateFieldsById(ctx, q.XId, bson.M{"status": "generating"})
	if err != nil {
		return 0, err
	}

	var answerText string

	var loop int64

	for {

		//if helper.InSlice(params.GetPrompt().GetId(), []string{"preAnalysisImages", "coverAnalysisImages", "limitAnalysisImages", "analysisImages"}) {
		//	model = genaiClient.GenerativeModel("gemini-2.5-pro-preview-05-06")
		//}

		//model.SystemInstruction = genai.NewUserContent(parts...)

		//model := genaiClient.GenerativeModel("gemini-1.5-pro-002")

		for resp, err := range genaiClient.GenerateContentStream(ctx, gemini.GenerateContentRequest{
			Model: modelName,
			Parts: parts,
		}) {
			if err != nil {
				//ERROR iterate fail err=googleapi: Error 403: You do not have permission to access the File xfr2ktwx3cb5 or it may not exist.
				if strings.Contains(err.Error(), "BlockReasonOther") {
					modelName = "gemini-2.5-pro-preview-05-06"
				}

				logger.Errorw("iterate fail err", err.Error())

				break
				//return 500, err
			}

			part := gemini.ResponseToString(resp)

			if strings.TrimSpace(part) == "" {
				continue
			}

			answerText += part
			esCtx.Write(part)
		}

		if answerText != "" {
			break
		}

		loop += 1
		logger.Debugw("answerText empty", "", "loop", loop)

		if loop > 3 {
			break
		}
	}

	if answerText == "" {
		return 500, errors.BadRequest("invalidAnswer", "")
	}

	// answer落库
	answerId := primitive.NewObjectID().Hex()
	newAnswer := &trackerpb.Answer{
		XId: answerId,
		Question: &trackerpb.Question{
			XId: q.XId,
		},
		Text:      answerText,
		UserId:    userId,
		CreatedAt: time.Now().Unix(),
	}
	_, err = t.Data.Mongo.Answer.InsertOne(ctx, newAnswer)
	if err != nil {
		return 500, err
	}

	_, err = t.Data.Mongo.Question.UpdateOne(ctx,
		bson.M{"_id": q.XId},
		bson.M{
			"$push": bson.M{"answers": newAnswer},
			"$set": bson.M{
				"status": "completed",
			},
		},
		options.Update().SetUpsert(true))
	if err != nil {
		return 500, err
	}

	question, err := t.Data.Mongo.Question.GetById(ctx, q.XId)
	if err != nil {
		return 500, err
	}

	esCtx.WriteDone(question.Safe())
	return 0, nil
}

func (t TrackerService) CreateQuestion(ctx context.Context, params *trackerpb.CreateQuestionParams) (int, error) {
	esCtx := ctx.(eventsource.Ctx)

	userId := esCtx.UserId()

	if !params.Retry {
		if !helper.InSlice(params.GetPrompt().GetId(), []string{"extractText", "editingTechAnalysis", "effectAnalysis", "soundAnalysis"}) {
			list, err := t.Data.Mongo.Question.List(ctx, bson.M{
				"userId":      userId,
				"session._id": params.SessionId,
				"prompt.id":   params.Prompt.GetId(),
				//"status":      bson.M{"$in": []string{ "generating", "created"}},
			})
			if err != nil {
				return 0, err
			}

			if len(list) != 0 {
				log.Debugw("question already exists", "", "userId", userId, "sessionId", params.SessionId, "promptId", params.Prompt.GetId())
				return 0, errors.BadRequest("duplicate", "")
			}
		}
	}

	genaiClient := t.Data.GenaiFactory.Get()

	session, err := t.Data.Mongo.Session.GetById(ctx, params.SessionId)
	if err != nil {
		return 0, errors.BadRequest("invalidSession", "")
	}

	if len(session.Resources) == 0 {
		return 0, errors.BadRequest("emptyResource", "")
	}

	//
	params.Prompt.SystemContent, err = t.GetPromptV2(ctx, params.Prompt, PromptArgs{
		AuthorProfile:   session.Resources[0].GetProfile().Text(),
		PersonalProfile: params.GetProfile().Text(),
	})
	if err != nil {
		return 0, err
	}

	if params.Prompt.SystemContent == "" {
		return 10400, errors.BadRequest("invalidPrompt", "")
	}

	log.Debugw("params", params, "prompt", params.Prompt)

	// 校验credit
	cost := configs.GetCost(params.Scene, params.Prompt.GetId())

	// 落库
	q, err := t.Data.Mongo.Question.InsertIfNX(ctx,
		bson.M{
			"session._id": params.SessionId,
			"prompt.id":   params.Prompt.GetId(),
			"status":      "created",
		},
		&trackerpb.Question{
			XId:       primitive.NewObjectID().Hex(),
			Session:   session,
			Prompt:    params.Prompt,
			UserId:    userId,
			CreatedAt: time.Now().Unix(),
			Cost: &trackerpb.Cost{
				//Id:        credit.CostId,
				Amount: cost,
				//Remaining: credit.Remaining,
			},
			Status: "created",
		})
	if err != nil {
		return 500, err
	}

	log.Debugw("InsertOne", "", "question.id", q.XId)
	//

	var parts []*genai.Part

	for _, x := range session.Resources {

		p, err := gemini.NewMediaPart(x.Url, x.MimeType)
		if err != nil {
			log.Errorw("NewMediaPart err", err, "url", x.Url)
			return 0, err
		}

		parts = append(parts, p)
	}

	parts = append(parts, gemini.NewTextPart(params.Prompt.SystemContent))

	credit, err := t.Data.GrpcClients.PaymentClient.CostCredit(ctx, &paymentpb.CheckCreditParams{
		UserId: userId,
		Cost:   cost,
	})

	if err != nil {
		log.Errorw("CostCredit err", err, "userId", userId)
		return 500, err
	}

	if !credit.Ok {
		return 10501, errors.BadRequest("exceeded", "")
	}

	//parts := []genai.Part{
	//	genai.FileData{
	//		MIMEType: "video/mp4",
	//		URI:      session.Resource.GenaiUri,
	//	},
	//	genai.Text(params.Prompt.SystemContent),
	//}
	//

	//model := service.Data.GenaiClient.GenerativeModel("gemini-2.0-flash")
	//model := t.Data.GenaiClient.GenerativeModel("gemini-2.5-pro-exp-03-25")
	//model := genaiClient.GenerativeModel("gemini-2.5-pro-preview-05-06")

	modelName := configs.GetModel(params.Scene, params.Prompt.GetId())

	//model := genaiClient.GenerativeModel(modelName)

	//model := genaiClient.GenerativeModel("gemini-1.5-pro-002")

	_, err = t.Data.Mongo.Question.UpdateFieldsById(ctx, q.XId, bson.M{"status": "generating"})
	if err != nil {
		return 0, err
	}

	var answerText string

	for resp, err := range genaiClient.GenerateContentStream(ctx, gemini.GenerateContentRequest{
		Model: modelName,
		Parts: parts,
	}) {
		// blocked: prompt: BlockReasonOther
		if err != nil {
			//ERROR iterate fail err=googleapi: Error 403: You do not have permission to access the File xfr2ktwx3cb5 or it may not exist.

			log.Errorw("iterate fail err", err)

			if strings.Contains(err.Error(), "BlockReasonOther") {

				t.Data.GrpcClients.PaymentClient.RollbackCredit(ctx, &paymentpb.RollbackCreditParams{
					UserId: userId,
					Amount: cost,
				})

				return 10502, errors.BadRequest("blockReasonOther", "")
			}

			return 500, err
		}

		part := gemini.ResponseToString(resp)

		if strings.TrimSpace(part) == "" {
			continue
		}

		log.Debugw("answer part", part)

		answerText += part
		esCtx.Write(part)
	}

	if answerText == "" {
		return 500, errors.BadRequest("invalidAnswer", "")
	}

	// answer落库
	answerId := primitive.NewObjectID().Hex()
	newAnswer := &trackerpb.Answer{
		XId: answerId,
		Question: &trackerpb.Question{
			XId: q.XId,
		},
		Text:      answerText,
		UserId:    userId,
		CreatedAt: time.Now().Unix(),
	}
	_, err = t.Data.Mongo.Answer.InsertOne(ctx, newAnswer)
	if err != nil {
		return 500, err
	}

	_, err = t.Data.Mongo.Question.UpdateOne(ctx,
		bson.M{"_id": q.XId},
		bson.M{
			"$push": bson.M{"answers": newAnswer},
			"$set": bson.M{
				"status": "completed",
			},
		},
		options.Update().SetUpsert(true))
	if err != nil {
		return 500, err
	}

	question, err := t.Data.Mongo.Question.GetById(ctx, q.XId)
	if err != nil {
		return 500, err
	}

	esCtx.WriteDone(question.Safe())

	return 0, nil
}

func (t TrackerService) GetQuestion(ctx context.Context, params *trackerpb.GetQuestionParams) (*trackerpb.Question, error) {
	//userId := krathelper.RequireUserId(ctx)

	question, err := t.Data.Mongo.Question.GetById(ctx, params.Id)
	if err != nil {
		log.Errorw("GetQuestion err", err, "id", params.Id)
		return nil, err
	}

	if len(question.Answers) == 0 {
		answer, _ := t.Data.LocalCache.Get("answer:" + question.XId)
		question.Answers = append(question.Answers, &trackerpb.Answer{
			Text: answer.(string),
		})
	}

	return question, nil
}

func (t TrackerService) CreateQuestionV2(ctx context.Context, params *trackerpb.CreateQuestionParams) (*trackerpb.Question, error) {

	userId := krathelper.RequireUserId(ctx)

	session, err := t.Data.Mongo.Session.GetById(ctx, params.SessionId)
	if err != nil {
		return nil, errors.BadRequest("invalidSession", "")
	}

	// todo  改成resourceId
	if len(session.Resources) == 0 {
		return nil, errors.BadRequest("emptyResource", "")
	}

	//
	params.Prompt.SystemContent, err = t.GetPromptV2(ctx, params.Prompt, PromptArgs{
		AuthorProfile:   session.Resources[0].GetProfile().Text(),
		PersonalProfile: params.GetProfile().Text(),
	})
	if err != nil {
		return nil, err
	}

	if params.Prompt.SystemContent == "" {
		return nil, errors.BadRequest("invalidPrompt", "")
	}

	// 校验credit
	credit, err := t.Data.GrpcClients.UserClient.CostCredit(ctx, &userpb.CheckCreditParams{
		UserId: userId,
		Cost:   10,
	})

	if err != nil {
		log.Errorw("CostCredit err", err, "userId", userId)
		return nil, err
	}

	if !credit.Ok {
		return nil, errors.BadRequest("exceeded", "")
	}

	log.Debugw("params", params, "prompt", params.Prompt)

	// 落库
	questionId := primitive.NewObjectID().Hex()

	newQuestion := &trackerpb.Question{
		XId:       questionId,
		Session:   session,
		Prompt:    params.Prompt,
		UserId:    userId,
		CreatedAt: time.Now().Unix(),
		Cost: &trackerpb.Cost{
			Id:        credit.CostId,
			Amount:    credit.Amount,
			Remaining: credit.Remaining,
		},
		Status: "prepared",
	}

	_, err = t.Data.Mongo.Question.InsertOne(ctx, newQuestion)
	if err != nil {
		return nil, err
	}

	return newQuestion, nil
}

//func (t TrackerService) AnswerQuestion(ctx context.Context) error {
//
//	list, err := t.Data.Mongo.Question.List(ctx, bson.M{"status": "prepared"})
//	if err != nil {
//		log.Errorw("ListQuestion err", err)
//		return err
//	}
//
//	if len(list) == 0 {
//		return nil
//	}
//
//	for _, x := range list {
//
//		log.Debugw("answer question", x)
//
//		_, err := t.Data.Mongo.Question.UpdateFieldsById(ctx, x.XId, bson.M{"status": "generating"})
//		if err != nil {
//			log.Errorw("UpdateFields err", err, "id", x.XId)
//			continue
//		}
//
//		go func(c context.Context, q *trackerpb.Question) {
//			defer helpers.DeferFunc()
//			err := t.answerQuestion(c, q)
//			if err != nil {
//				log.Errorw("answer question err", err, "id", x.XId)
//				return
//			}
//		}(ctx, x)
//	}
//
//	return nil
//}
//
//func (t TrackerService) answerQuestion(ctx context.Context, question *trackerpb.Question) error {
//
//	genaiClient := t.Data.GenaiFactory.Get()
//
//	var parts []genai.Part
//
//	for _, x := range question.GetSession().Resources {
//		parts = append(parts, genai.FileData{
//			MIMEType: x.MimeType,
//			URI:      x.GenaiUri,
//		})
//	}
//
//	parts = append(parts, genai.Text(question.Prompt.SystemContent))
//
//	//parts := []genai.Part{
//	//	genai.FileData{
//	//		MIMEType: "video/mp4",
//	//		URI:      session.Resource.GenaiUri,
//	//	},
//	//	genai.Text(params.Prompt.SystemContent),
//	//}
//	//
//
//	//model := service.Data.GenaiClient.GenerativeModel("gemini-2.0-flash")
//	//model := t.Data.GenaiClient.GenerativeModel("gemini-2.5-pro-exp-03-25")
//	model := genaiClient.GenerativeModel("gemini-2.5-pro-preview-05-06")
//
//	iter := model.GenerateContentStream(ctx, parts...)
//	//resp, err := model.GenerateContent(ctx, parts...)
//	//if err != nil {
//	//	log.Error(err)
//	//	ectx.Abort(500, err.Error())
//	//	return
//	//}
//	//
//	//answer := service.ResponseString(resp)
//
//	var answerText string
//
//	for {
//		resp, err := iter.Next()
//		if errors.Is(err, iterator.Done) {
//			break
//		}
//
//		if err != nil {
//
//			var e *apierror.APIError
//			ok := errors.As(err, &e)
//			if ok {
//				log.Errorw("iterate fail err", e.Unwrap())
//			} else {
//				log.Errorw("iterate fail err", err)
//			}
//
//		}
//
//		part := t.ResponseString(resp)
//
//		log.Debugw("answer part", part)
//
//		answerText += part
//
//		t.Data.LocalCache.Set("answer:"+question.XId, answerText, -1)
//	}
//
//	// answer落库
//	answerId := primitive.NewObjectID().Hex()
//	newAnswer := &trackerpb.Answer{
//		XId: answerId,
//		Question: &trackerpb.Question{
//			XId: question.XId,
//		},
//		Text:      answerText,
//		UserId:    question.UserId,
//		CreatedAt: time.Now().Unix(),
//	}
//	_, err := t.Data.Mongo.Answer.InsertOne(ctx, newAnswer)
//	if err != nil {
//		return err
//	}
//
//	_, err = t.Data.Mongo.Question.UpdateOne(ctx,
//		bson.M{"_id": question.XId},
//		bson.M{
//			"$push": bson.M{"answers": newAnswer},
//			"$set": bson.M{
//				"status": "completed",
//			},
//		},
//		options.UpdateFields().SetUpsert(true))
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
