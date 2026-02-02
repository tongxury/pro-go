package events

import (
	"store/pkg/sdk/conv"
	"time"
)

type SubscribeEvent struct {
	ID       string
	UserID   string
	Level    string
	Cycle    string
	Renew    bool
	Ts       int64
	Manually bool
}

type AuthEvent struct {
	UserID     string
	LoginBy    string
	TS         int64
	DeviceID   string
	IsRegister bool
}

type AuthEventValues map[string]interface{}

func NewAuthEvent(userID string, loginBy, deviceID string) map[string]interface{} {
	t := &AuthEvent{
		UserID:   userID,
		LoginBy:  loginBy,
		TS:       time.Now().Unix(),
		DeviceID: deviceID,
	}
	return t.AsValues()
}

func (t *AuthEvent) AsValues() map[string]interface{} {
	return map[string]any{
		"userID":   t.UserID,
		"loginBy":  t.LoginBy,
		"ts":       time.Now().Unix(),
		"deviceID": t.DeviceID,
	}
}

func (t AuthEventValues) AsAuthEvent() *AuthEvent {
	return &AuthEvent{
		UserID:   conv.String(t["userID"]),
		LoginBy:  conv.String(t["loginBy"]),
		TS:       conv.Int64(t["ts"]),
		DeviceID: conv.String(t["deviceID"]),
	}
}
