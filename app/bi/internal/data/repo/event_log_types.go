package repo

import "time"

type EventLog struct {
	EventName   string
	Ip          string
	CountryCode string
	CreatedAt   time.Time
	UserID      string
	DeviceID    string
	Channel     string
	Platform    string
	Idfa        string
	Imei        string
	Os          string
	Oaid        string
}

type EventLogs []EventLog

type EventStats struct {
	Visit     uint64 `ch:"visit"`
	Install   uint64 `ch:"install"`
	Register  uint64 `ch:"register"`
	Subscribe uint64 `ch:"subscribe"`
}

type StatItem struct {
	EventID     string `ch:"event_id"`
	DeviceCount uint64 `ch:"device_count"`
	UserCount   uint64 `ch:"user_count"`
}

type StatItems []StatItem

func (ts StatItems) AsMap() map[string]StatItem {
	rsp := make(map[string]StatItem, len(ts))

	for _, t := range ts {
		rsp[t.EventID] = t
	}

	return rsp
}
