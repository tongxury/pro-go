package repo

import (
	"context"
	"store/pkg/clients"
	"time"
)

type ChatRecordRepo struct {
	db *clients.ClickHouseClient
}

func NewChatRecordRepo(db *clients.ClickHouseClient) *ChatRecordRepo {
	return &ChatRecordRepo{db: db}
}

func (t *ChatRecordRepo) AsyncInsert(ctx context.Context, e *ChatRecord) error {

	err := t.db.AsyncInsert(ctx, `
		insert into chat_records (event_time, user_id, device_id, function_name, model, url, query, image, answer, status)
				values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, false,
		e.EventTime.Format(time.DateTime), e.UserID, e.DeviceID, e.FunctionName, e.Model,
		e.Url, e.Query, e.Image, e.Answer, e.Status,
	)
	if err != nil {
		return err
	}

	return nil
}

type ListChatRecordParams struct {
	Page, Size int64
	UserIDs    []string
}

func (t *ChatRecordRepo) ListChatRecords(ctx context.Context, params *ListChatRecordParams) error {

	return nil
}
