package repo

import (
	"context"
	"store/pkg/clients"
	"store/pkg/sdk/helper"
	"time"
)

type EventLogRepo struct {
	db *clients.ClickHouseClient
}

func NewEventLogRepo(db *clients.ClickHouseClient) *EventLogRepo {
	return &EventLogRepo{db: db}
}

func (t *EventLogRepo) Insert(ctx context.Context, e EventLog) error {

	err := t.db.Exec(ctx, `
				insert into events (
                        event_name,
                        ip, 
                        country_code, 
                        created_at, 
                        user_id, 
                        device_id,
				        channel,
				        platform,
				        idfa,
						imei,
						oaid
				        ) 
				values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.EventName, helper.OrString(e.Ip, "0.0.0.0"), e.CountryCode, e.CreatedAt.Format(time.DateTime),
		e.UserID, e.DeviceID, e.Channel, e.Platform, e.Idfa, e.Imei, e.Oaid,
	)
	if err != nil {
		return err
	}

	return nil
}

func (t *EventLogRepo) AsyncInsert(ctx context.Context, e EventLog) error {

	err := t.db.AsyncInsert(ctx, `
				insert into events (
                        event_name,
                        ip, 
                        country_code, 
                        created_at, 
                        user_id, 
                        device_id,
				        channel,
				        platform,
				        idfa,
						imei,
						oaid
				        ) 
				values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, false,
		e.EventName, helper.OrString(e.Ip, "0.0.0.0"), e.CountryCode, e.CreatedAt.Format(time.DateTime),
		e.UserID, e.DeviceID, e.Channel, e.Platform, e.Idfa, e.Imei, e.Oaid,
	)
	if err != nil {
		return err
	}

	return nil
}

func (t *EventLogRepo) InsertInBulk(ctx context.Context, es EventLogs) error {

	batch, err := t.db.PrepareBatch(ctx, `
			insert into event_logs (
                        event_name,
                        ip, 
                        country_code, 
                        created_at, 
                        user_id, 
                        device_id,
				        channel,
			                        platform
				        ) 
	`)
	if err != nil {
		return err
	}

	for _, e := range es {
		err := batch.Append(
			e.EventName, e.Ip, e.CountryCode, e.CreatedAt.Format(time.DateTime),
			e.UserID, e.DeviceID, e.Channel, e.Platform,
		)
		if err != nil {
			return err
		}
	}

	return batch.Send()
}
