package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/krathelper"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *VoiceAgentService) CreateEvent(ctx context.Context, req *voiceagent.CreateEventRequest) (*voiceagent.ImportantEvent, error) {
	userId := krathelper.RequireUserId(ctx)

	event := &voiceagent.ImportantEvent{
		XId:           primitive.NewObjectID().Hex(),
		UserId:        userId,
		Title:         req.Title,
		Type:          req.Type,
		Date:          req.Date,
		IsRecurring:   req.IsRecurring,
		Note:          req.Note,
		RelatedPerson: req.RelatedPerson,
		ReminderDays:  req.ReminderDays,
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
	}

	// 默认提前1天提醒
	if event.ReminderDays == 0 {
		event.ReminderDays = 1
	}

	res, err := s.Data.Mongo.ImportantEvent.Insert(ctx, event)
	if err != nil {
		return nil, err
	}

	return res, nil
}
