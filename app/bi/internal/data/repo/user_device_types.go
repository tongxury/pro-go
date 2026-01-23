package repo

import "time"

type UserDevice struct {
	DeviceID      string    `ch:"device_id"`
	Ip            string    `ch:"ip"`
	AppVersion    string    `ch:"app_version"`
	CreatedAt     time.Time `ch:"created_at"`
	ChromeVersion string    `ch:"chrome_version"`
	Platform      string    `ch:"platform"`
	UserID        int64     `ch:"user_id"`
	Channel       string    `ch:"channel"`
	Country       string    `ch:"country"`
}

type UserDevices []UserDevice

func (ts UserDevices) AsMap() map[int64]UserDevice {

	rsp := make(map[int64]UserDevice, len(ts))

	for _, t := range ts {
		rsp[t.UserID] = t
	}

	return rsp
}
