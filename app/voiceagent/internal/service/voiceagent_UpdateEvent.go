package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"time"
)

func (s *VoiceAgentService) UpdateEvent(ctx context.Context, req *voiceagent.UpdateEventRequest) (*voiceagent.ImportantEvent, error) {
	userId := krathelper.RequireUserId(ctx)

	// 确保只能更新自己的事件
	event, err := s.Data.Mongo.ImportantEvent.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if event.UserId != userId {
		return nil, krathelper.ErrForbidden
	}

	// 构建更新操作
	updateOp := mgz.Op().Set("updatedAt", time.Now().Unix())

	if req.Title != "" {
		updateOp = updateOp.Set("title", req.Title)
	}
	if req.Type != "" {
		updateOp = updateOp.Set("type", req.Type)
	}
	if req.Date != "" {
		updateOp = updateOp.Set("date", req.Date)
	}
	updateOp = updateOp.Set("isRecurring", req.IsRecurring)
	if req.Note != "" {
		updateOp = updateOp.Set("note", req.Note)
	}
	if req.RelatedPerson != "" {
		updateOp = updateOp.Set("relatedPerson", req.RelatedPerson)
	}
	if req.ReminderDays > 0 {
		updateOp = updateOp.Set("reminderDays", req.ReminderDays)
	}

	_, err = s.Data.Mongo.ImportantEvent.UpdateByIDIfExists(ctx, req.Id, updateOp)
	if err != nil {
		return nil, err
	}

	return s.Data.Mongo.ImportantEvent.GetById(ctx, req.Id)
}
