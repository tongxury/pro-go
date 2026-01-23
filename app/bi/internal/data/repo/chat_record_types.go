package repo

import "time"

type ChatRecord struct {
	EventTime    time.Time `ch:"event_time"`
	UserID       int64     `ch:"user_id"`
	DeviceID     string    `ch:"device_id"`
	FunctionName string    `ch:"function_name"`
	Model        string    `ch:"model"`
	Url          string    `ch:"url"`
	Query        string    `ch:"query"`
	Image        string    `ch:"image"`
	Answer       string    `ch:"answer"`
	Status       string    `ch:"status"`
}
