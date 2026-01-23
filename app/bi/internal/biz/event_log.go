package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"store/app/bi/internal/data"
	"store/app/bi/internal/data/repo"
	"time"
)

type EventLogBiz struct {
	data *data.Data
}

func NewEventLogBiz(data *data.Data) *EventLogBiz {
	return &EventLogBiz{data: data}
}

type AddEventLogParams struct {
	EventName   string
	IP          string
	CountryCode string
	CreatedAt   time.Time
	UserID      string
	DeviceID    string
	Channel     string
	Platform    string
	Idfa        string
	Imei        string
	Oaid        string
}

func (t *EventLogBiz) AddEventLog(ctx context.Context, params AddEventLogParams) error {

	log.Debug("AddEventLog ing", params)

	err := t.data.Repos.EventLog.AsyncInsert(ctx, repo.EventLog{
		EventName:   params.EventName,
		Ip:          params.IP,
		CountryCode: params.CountryCode,
		CreatedAt:   params.CreatedAt,
		UserID:      params.UserID,
		DeviceID:    params.DeviceID,
		Channel:     params.Channel,
		Platform:    params.Platform,
		Idfa:        params.Idfa,
		Imei:        params.Imei,
		Oaid:        params.Oaid,
	})
	if err != nil {
		return err
	}

	return nil
}
