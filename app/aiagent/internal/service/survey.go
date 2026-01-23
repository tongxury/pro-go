package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/protobuf/ptypes/empty"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/genai"
	aiagentpb "store/api/aiagent"
	"store/pkg/krathelper"
	helpers "store/pkg/sdk/helper"
	"store/pkg/sdk/third/gemini"
	"strings"
)

func (t TrackerService) UpdateSurvey(ctx context.Context, params *aiagentpb.UpdateSurveyParams) (*empty.Empty, error) {

	_, err := t.Data.Mongo.Survey.UpdateFieldsById(ctx, params.Id, bson.M{"status": "unlocked"})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &empty.Empty{}, nil

}

func (t TrackerService) GetSurvey(ctx context.Context, params *aiagentpb.GetSurveyParams) (*aiagentpb.Survey, error) {

	survey, err := t.Data.Mongo.Survey.GetById(ctx, params.Id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if survey.Status != "unlocked" {
		survey.Result = ""
	}

	return survey, nil
}

func (t TrackerService) GenerateSurveyResult(ctx context.Context) error {

	list, err := t.Data.Mongo.Survey.List(ctx, bson.M{"status": "submitted"})
	if err != nil {
		log.Errorw("List err", err)
		return err
	}

	if len(list) == 0 {
		return nil
	}

	for _, x := range list {

		log.Debugw("GenerateSurveyResult", x)

		_, err := t.Data.Mongo.Survey.UpdateFieldsById(ctx, x.XId, bson.M{"status": "generating"})
		if err != nil {
			log.Errorw("UpdateFields err", err, "id", x.XId)
			continue
		}

		go func(c context.Context, q *aiagentpb.Survey) {
			defer helpers.DeferFunc()
			err := t.generate(c, q)
			if err != nil {
				log.Errorw("GenerateSurveyResult err", err, "id", x.XId)
				return
			}
		}(ctx, x)
	}

	return nil
}

func (t TrackerService) generate(ctx context.Context, survey *aiagentpb.Survey) error {
	//var parts []genai.Part

	genaiClient := t.Data.GenaiFactory.Get()

	settings, err := t.Data.Mongo.Settings.Get(ctx)
	if err != nil {
		return err
	}

	prompt := settings.Prompt["survey"]

	prompt = strings.ReplaceAll(prompt, "__SURVEY_FORM__", survey.Text)

	//parts = append(parts, genai.Text(survey.Text))
	parts := []*genai.Part{
		gemini.NewTextPart(survey.Text),
		gemini.NewTextPart(prompt),
	}

	log.Debugw("survey prompt", prompt)

	//model := service.Data.GenaiClient.GenerativeModel("gemini-2.0-flash")
	//model := t.Data.GenaiClient.GenerativeModel("gemini-2.5-pro-exp-03-25")
	modelName := "gemini-2.5-pro-preview-05-06"

	answer, err := genaiClient.GenerateContent(ctx, gemini.GenerateContentRequest{
		Model: modelName,
		Parts: parts,
	})
	if err != nil {
		log.Error(err)
		return nil
	}

	_, err = t.Data.Mongo.Survey.UpdateFieldsById(ctx, survey.XId, bson.M{
		"result": answer, "status": "generated",
	})
	if err != nil {
		return err
	}

	return nil
}

func (t TrackerService) SubmitSurvey(ctx context.Context, params *aiagentpb.SubmitSurveyParams) (*aiagentpb.Survey, error) {

	deviceId := krathelper.GetHeader(ctx, "Device-Id")

	newSurvey := &aiagentpb.Survey{
		XId:      primitive.NewObjectID().Hex(),
		Text:     params.Text,
		DeviceId: deviceId,
		Status:   "submitted",
	}

	err := t.Data.Mongo.Survey.Insert(ctx, newSurvey)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return newSurvey, nil
}
