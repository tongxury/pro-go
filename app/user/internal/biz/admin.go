package biz

import (
	"context"
	typepb "store/api/admin/types"
	"store/app/user/internal/data"
	"store/app/user/internal/data/repo/ent/notification"
	"store/pkg/enums"
	"store/pkg/sdk/conv"
	"time"
)

type AdminBiz struct {
	data *data.Data
}

func NewAdminBiz(data *data.Data) *AdminBiz {
	return &AdminBiz{data: data}
}

func (t *AdminBiz) GetNotification(ctx context.Context, userId string) (*typepb.Notification, error) {

	now := time.Now()

	items, err := t.data.Repos.EntClient.Notification.Query().
		Where(notification.StartAtLT(now)).
		Where(notification.EndAtGT(now)).
		All(ctx)

	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}

	item := items[0]

	var contents []*typepb.Notification_Content
	for _, x := range item.Extra.Contents {
		contents = append(contents, &typepb.Notification_Content{
			Type:  x.Type,
			Title: x.Title,
			Value: x.Value,
		})
	}

	var id string
	if item.Level == enums.NotificationLevel_Info {
		id = conv.String(item.ID)
	}

	return &typepb.Notification{
		Id:       id,
		Level:    item.Level.String(),
		Title:    item.Extra.Title,
		Contents: contents,
	}, nil
}
